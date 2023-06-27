package matcher

import (
	"sync"

	"github.com/asaskevich/EventBus"
	"github.com/ksysoev/gochess/events"
)

type Matcher struct {
	EventBus EventBus.Bus
	queue    []string
	lock     sync.Mutex
}

func NewMatcher(evbus EventBus.Bus) Matcher {
	return Matcher{
		EventBus: evbus,
		queue:    []string{},
		lock:     sync.Mutex{},
	}
}

func (m *Matcher) findMatch(player string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.queue) == 0 {
		m.queue = append(m.queue, player)
		return nil
	}

	white := m.queue[0]
	m.queue = m.queue[1:]
	black := player

	m.EventBus.Publish("match::found", events.MatchFoundEvent{
		White: white,
		Black: black,
	})
	return nil
}
