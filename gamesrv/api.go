package gamesrv

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ApiGameServer struct {
	GameService GameService
	Router      chi.Router
}

func NewApiGameServer() ApiGameServer {
	app := ApiGameServer{
		GameService: GameService{
			GameRepo: NewGameRepo(),
		},
	}

	r := chi.NewRouter()
	r.Post("/", app.CreateGame)
	r.Get("/{gameID}", app.GetGame)
	r.Post("/{gameID}/move", app.MakeMove)

	app.Router = r

	return app
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

func (app *ApiGameServer) CreateGame(w http.ResponseWriter, r *http.Request) {
	// Start a new game

	game, err := app.GameService.CreateGame("player1", "player2")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := StartGameResponse{
		ID:       game.ID,
		Position: game.Position,
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

func (app *ApiGameServer) GetGame(w http.ResponseWriter, r *http.Request) {

	gameID := chi.URLParam(r, "gameID")
	if gameID == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	game, err := app.GameService.GetGame(gameID)

	if err != nil {
		http.Error(w, "Game not started or finished", http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: game.Position,
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

func (app *ApiGameServer) MakeMove(w http.ResponseWriter, r *http.Request) {

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

	game, err := app.GameService.MakeMove(gameID, req.Move)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := MoveResponse{
		Position: game.Position,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Illigal move", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
