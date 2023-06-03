package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/notnil/chess"
)

type MoveRequest struct {
	Move string `json:"move"`
}

type MoveResponse struct {
	Position string `json:"position"`
}

var currentGame *chess.Game

func startGame(w http.ResponseWriter, r *http.Request) {
	// Start a new game
	currentGame = chess.NewGame()

	w.WriteHeader(http.StatusOK)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Error(w, "Game not started", http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: currentGame.Position().String(),
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
	if currentGame == nil {
		http.Error(w, "Game not started", http.StatusBadRequest)
		return
	}

	currentGame = nil

	w.WriteHeader(http.StatusOK)
}

func move(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Error(w, "Game not started", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := currentGame.MoveStr(req.Move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: currentGame.Position().String(),
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
	mux := http.NewServeMux()

	// Register the /start, /move, and /finish routes
	mux.HandleFunc("/start", startGame)
	mux.HandleFunc("/move", move)
	mux.HandleFunc("/finish", finishGame)
	mux.HandleFunc("/game", getGame)

	// Serve the routes using the ServeMux
	log.Println("Starting Game Server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
