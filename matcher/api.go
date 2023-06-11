package matcher

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MatcherAPI struct {
	Routes  chi.Router
	Service Matcher
}

func NewMatcherAPI(bus EventBus.Bus) MatcherAPI {
	matcher := NewMatcher(bus)

	api := MatcherAPI{
		Service: matcher,
	}

	r := chi.NewRouter()
	r.Post("/", api.findMatch)

	api.Routes = r

	return api
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

func (api *MatcherAPI) findMatch(w http.ResponseWriter, r *http.Request) {
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
		// TODO: create game ... probably we need service discovery

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
