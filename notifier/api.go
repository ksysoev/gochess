package notifier

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi"
)

// Event is a struct that represents a single event
type Event struct {
	ID   string
	Data string
	Name string
}

// ApiNotifierServer is a struct that represents a notifier server
type ApiNotifierServer struct {
	Router   chi.Router
	Notifier NotifierService
}

// NewApiNotifierServer is a function that creates a new notifier server
func NewApiNotifierServer(evbus EventBus.Bus) ApiNotifierServer {
	app := ApiNotifierServer{
		Notifier: NewNotifierService(evbus),
	}

	r := chi.NewRouter()
	r.Get("/notifier", app.SubscribeEvents)

	app.Router = r

	return app
}

// SubscribeEvents is a function that subscribes to events
func (app *ApiNotifierServer) SubscribeEvents(w http.ResponseWriter, r *http.Request) {
	flusher, err := w.(http.Flusher)
	if !err {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	eventChan := app.Notifier.Subscribe([]string{"game:move", "game:start"})

	flusher.Flush()
	// Continuously read events from the channel and send them to the SSE stream
	for {
		timeout := time.After(15 * time.Second)
		select {
		case event := <-eventChan:
			if event.Name != "" {
				fmt.Fprintf(w, "event: %s\n", event.Name)
			}

			if event.ID != "" {
				fmt.Fprintf(w, "id: %s\n", event.ID)
			}

			if event.Data != "" {
				fmt.Fprintf(w, "data: %s\n", event.Data)
			}

			fmt.Fprintf(w, "\n")
		case <-timeout:
			fmt.Fprintf(w, "event: ping\n\n")
		}
		flusher.Flush()
	}
}
