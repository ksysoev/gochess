package matcher

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi/v5"
)

type MatcherAPI struct {
	Router  chi.Router
	Service Matcher
}

func NewMatcherAPI(bus EventBus.Bus) MatcherAPI {
	matcher := NewMatcher(bus)

	api := MatcherAPI{
		Service: matcher,
	}

	r := chi.NewRouter()
	r.Post("/", api.findMatch)

	api.Router = r

	return api
}

type MatchRequest struct {
	Name string `json:"name"`
}

type MatchResponse struct {
	Status string `json:"status"`
}

func (api *MatcherAPI) findMatch(w http.ResponseWriter, r *http.Request) {
	var req MatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	err = api.Service.findMatch(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MatchResponse{
		Status: "ok",
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
