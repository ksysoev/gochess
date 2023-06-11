package matcher

import "github.com/asaskevich/EventBus"

type Matcher struct {
	EventBus EventBus.Bus
}

func NewMatcher(evbus EventBus.Bus) Matcher {
	return Matcher{
		EventBus: evbus,
	}
}
