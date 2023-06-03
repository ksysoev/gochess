package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartGame(t *testing.T) {
	// Create a new request to the /start endpoint
	req, err := http.NewRequest("GET", "/start", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the startGame handler function
	handler := http.HandlerFunc(startGame)
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the current game is not nil and has started
	assert.NotNil(t, currentGame)
	assert.True(t, currentGame.Started)
}

func TestFinishGame(t *testing.T) {
	// Create a new request to the /finish endpoint
	req, err := http.NewRequest("GET", "/finish", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the finishGame handler function
	handler := http.HandlerFunc(finishGame)
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Start a new game
	currentGame = &Game{
		Board:    chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation)).Position().Board(),
		Players:  [2]string{"Player 1", "Player 2"},
		Turn:     chess.White,
		Started:  true,
		Finished: false,
	}

	// Call the finishGame handler function again
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the current game has finished
	assert.True(t, currentGame.Finished)
}

func TestMove(t *testing.T) {
	// Create a new request to the /move endpoint
	moveReq := MoveRequest{
		From: "e2",
		To:   "e4",
	}
	moveReqBytes, _ := json.Marshal(moveReq)
	req, err := http.NewRequest("POST", "/move", bytes.NewBuffer(moveReqBytes))
	if err != nil {
		t.Fatal(err)
	}

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Start a new game
	currentGame = &Game{
		Board:    chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation)).Position().Board(),
		Players:  [2]string{"Player 1", "Player 2"},
		Turn:     chess.White,
		Started:  true,
		Finished: false,
	}

	// Call the move handler function
	handler := http.HandlerFunc(move)
	handler.ServeHTTP(rr, req)

	// Check that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the current game state has been updated
	var gameState GameState
	json.NewDecoder(rr.Body).Decode(&gameState)
	assert.NotNil(t, gameState.Board)
	assert.Equal(t, "e4", gameState.Board.Squares()[12].Name())
	assert.Equal(t, "w", gameState.Turn)
}
