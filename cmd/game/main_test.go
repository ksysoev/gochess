package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/notnil/chess"
)

func TestStartGame(t *testing.T) {
	// Create a new HTTP request to the /game route
	req, err := http.NewRequest("POST", "/game", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Create a new chi router and register the /game route
	r := chi.NewRouter()
	r.Post("/game", startGame)

	// Send the HTTP request to the chi router
	r.ServeHTTP(rr, req)

	// Check the status code of the HTTP response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type of the HTTP response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, "application/json")
	}

	// Check the body of the HTTP response
	var resp StartGameResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if resp.ID == "" {
		t.Errorf("handler returned empty ID")
	}

	if resp.Position == "" {
		t.Errorf("handler returned empty position")
	}

	// Check that the game position was saved to the repository
	if pos, ok := GameRepo.Postitions[resp.ID]; !ok || pos != resp.Position {
		t.Errorf("handler did not save game position to repository")
	}
}

func TestGetGame(t *testing.T) {
	// Create a new game and save its position to the repository
	newGame := chess.NewGame()
	id := uuid.New().String()
	GameRepo.Postitions[id] = newGame.Position().String()

	// Create a new HTTP request to the /game/{gameID} route
	req, err := http.NewRequest("GET", fmt.Sprintf("/game/%s", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Create a new chi router and register the /game/{gameID} route
	r := chi.NewRouter()
	r.Get("/game/{gameID}", getGame)

	// Send the HTTP request to the chi router
	r.ServeHTTP(rr, req)

	// Check the status code of the HTTP response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type of the HTTP response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, "application/json")
	}

	// Check the body of the HTTP response
	var resp MoveResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if resp.Position == "" {
		t.Errorf("handler returned empty position")
	}

	// Check that the game position in the response matches the one in the repository
	if resp.Position != GameRepo.Postitions[id] {
		t.Errorf("handler returned incorrect game position")
	}
}

func TestMove(t *testing.T) {
	// Create a new game and save its position to the repository
	newGame := chess.NewGame()
	id := uuid.New().String()
	GameRepo.Postitions[id] = newGame.Position().String()

	// Create a new HTTP request to the /game/{gameID}/move route
	moveReq := MoveRequest{
		Move: "e4",
	}
	moveReqBytes, err := json.Marshal(moveReq)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("/game/%s/move", id), bytes.NewBuffer(moveReqBytes))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Create a new chi router and register the /game/{gameID}/move route
	r := chi.NewRouter()
	r.Post("/game/{gameID}/move", move)

	// Send the HTTP request to the chi router
	r.ServeHTTP(rr, req)

	// Check the status code of the HTTP response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type of the HTTP response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			ctype, "application/json")
	}

	// Check the body of the HTTP response
	var resp MoveResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if resp.Position == "" {
		t.Errorf("handler returned empty position")
	}

	// Check that the game position in the response matches the one in the repository
	if resp.Position != GameRepo.Postitions[id] {
		t.Errorf("handler returned incorrect game position")
	}
}
