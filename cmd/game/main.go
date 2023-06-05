package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/notnil/chess"

	"github.com/ksysoev/gochess/gamesrv"
)

type GameRepoStorage struct {
	Postitions map[string]string
}

var GameRepo GameRepoStorage = GameRepoStorage{
	Postitions: make(map[string]string),
}

var queue []string = make([]string, 0)

type MatchRequest struct {
	Name string `json:"name"`
}

type MatchResponse struct {
	White  string `json:"white,omitempty"`
	Black  string `json:"black,omitempty"`
	Status string `json:"status"`
	GameID string `json:"game_id,omitempty"`
}

func findMatch(w http.ResponseWriter, r *http.Request) {
	var req MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var resp MatchResponse
	if len(queue) == 0 {
		queue = append(queue, req.Name)
		resp = MatchResponse{
			Status: "pending",
		}
	} else {
		white := queue[0]
		queue = queue[1:]
		black := req.Name
		id := uuid.New().String()
		newGame := chess.NewGame()
		GameRepo.Postitions[id] = newGame.Position().String()

		resp = MatchResponse{
			Status: "ready",
			White:  white,
			Black:  black,
			GameID: id,
		}
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

	// CORS for development purposes
	r.Use(SetAccessControlHeader)

	r.Post("/match", findMatch)

	gamesrv := gamesrv.NewApiGameServer()
	r.Mount("/game", gamesrv.Router)

	// Serve the routes using the ServeMux
	log.Println("Starting Game Server on port 8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}

func SetAccessControlHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
