package notifier

import "github.com/asaskevich/EventBus"

type NotifierService struct {
	EventBus EventBus.Bus
}

func NewNotifierService(evbus EventBus.Bus) NotifierService {
	return NotifierService{
		EventBus: evbus,
	}
}
