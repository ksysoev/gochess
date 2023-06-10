package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ksysoev/gochess/gamesrv"
)

type Event struct {
	ID   string
	Data string
	Name string
}

type ApiNotifierServer struct {
	Router   chi.Router
	Notifier NotifierService
}

func NewApiNotifierServer(evbus EventBus.Bus) ApiNotifierServer {
	app := ApiNotifierServer{
		Notifier: NewNotifierService(evbus),
	}

	r := chi.NewRouter()
	r.Get("/notifier", app.SubscribeEvents)

	app.Router = r

	return app
}

func (app *ApiNotifierServer) SubscribeEvents(w http.ResponseWriter, r *http.Request) {
	flusher, err := w.(http.Flusher)
	if !err {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	eventChan := make(chan Event)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe to all events on the bus and send them to the channel
	app.Notifier.EventBus.Subscribe("game:move", func(event interface{}) {
		eventObj, ok := event.(gamesrv.EventGameMove)

		if !ok {
			http.Error(w, "Invalid event type", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(eventObj)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eventChan <- Event{
			ID:   uuid.New().String(),
			Data: string(data),
			Name: "game:move",
		}
	})

	app.Notifier.EventBus.Subscribe("game:start", func(event interface{}) {
		eventObj, ok := event.(gamesrv.EventGameStart)

		if !ok {
			http.Error(w, "Invalid event type", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(eventObj)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("event: ", string(data))

		eventChan <- Event{
			ID:   uuid.New().String(),
			Data: string(data),
			Name: "game:start",
		}
	})

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
