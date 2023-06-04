package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/notnil/chess"
)

type GameRepoStorage struct {
	CurrentGameID int
	Postitions    map[string]string
}

var GameRepo GameRepoStorage = GameRepoStorage{
	CurrentGameID: 0,
	Postitions:    make(map[string]string),
}

type MoveRequest struct {
	Move string `json:"move"`
}

type MoveResponse struct {
	Position string `json:"position"`
}

type StartGameResponse struct {
	ID       string `json:"id"`
	Position string `json:"position"`
}

func startGame(w http.ResponseWriter, r *http.Request) {
	// Start a new game
	newGame := chess.NewGame()

	GameRepo.CurrentGameID += 1
	id := fmt.Sprintf("%d", GameRepo.CurrentGameID)
	GameRepo.Postitions[id] = newGame.Position().String()

	resp := StartGameResponse{
		ID:       id,
		Position: newGame.Position().String(),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	gamePosition, ok := GameRepo.Postitions[gameID]

	if !ok {
		http.Error(w, "Game not started or finished", http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: gamePosition,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func finishGame(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	delete(GameRepo.Postitions, gameID)

	w.WriteHeader(http.StatusOK)
}

func move(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameID")

	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gamePosition, ok := GameRepo.Postitions[gameID]

	if !ok {
		http.Error(w, "Game not started or finished", http.StatusBadRequest)
		return
	}

	fen, err := chess.FEN(gamePosition)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game := chess.NewGame(fen)

	if err := game.MoveStr(req.Move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: game.Position().String(),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	// Create a new ServeMux
	r := chi.NewRouter()

	// Register the /start, /move, and /finish routes
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Post("/game", startGame)
	r.Get("/game/{gameID}", getGame)
	r.Delete("/game/{gameID}", finishGame)
	r.Post("/game/{gameID}/move", move)

	// Serve the routes using the ServeMux
	log.Println("Starting Game Server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
