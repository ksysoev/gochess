package main

import (
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ksysoev/gochess/gamesrv"
	"github.com/ksysoev/gochess/matcher"
	"github.com/ksysoev/gochess/notifier"
)

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

	bus := EventBus.New()
	gamesrv := gamesrv.NewApiGameServer(bus)
	r.Mount("/game", gamesrv.Router)

	matcher := matcher.NewMatcherAPI(bus)
	r.Mount("/match", matcher.Router)

	srv := notifier.NewApiNotifierServer(bus)
	r.Mount("/notifier", srv.Router)
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
