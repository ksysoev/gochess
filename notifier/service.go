package notifier

import (
	"encoding/json"
	"log"

	"github.com/asaskevich/EventBus"
	"github.com/google/uuid"
)

// NotifierService is a service for notifying about events.
type NotifierService struct {
	EventBus EventBus.Bus
}

// NewNotifierService creates a new NotifierService.
func NewNotifierService(evbus EventBus.Bus) NotifierService {
	return NotifierService{
		EventBus: evbus,
	}
}

func (ns NotifierService) Subscribe(events []string) chan Event {
	eventChan := make(chan Event)

	for _, eventName := range events {
		ns.EventBus.Subscribe(eventName, func(event any) {
			data, err := json.Marshal(event)

			if err != nil {
				log.Println("Error in serializing event: ", err)
				return
			}

			eventChan <- Event{
				ID:   uuid.New().String(),
				Data: string(data),
				Name: eventName,
			}
		})
	}

	return eventChan
}
