package observability

import (
	"sync"
	"time"
)

type Event struct {
	TraceID   string    `json:"trace_id"`
	EventType string    `json:"event_type"`
	Outcome   string    `json:"outcome"`
	At        time.Time `json:"at"`
}

type Sink struct {
	// Guards concurrent slice mutations. Sink is a shared state component.
	mu     sync.Mutex
	events []Event
	// Strict bound to prevent memory exhaustion (OOM) during event spikes.
	limit int
}

func NewSink(limit int) *Sink {
	if limit < 1 {
		// Fallback boundary to prevent accidental zero-capacity deadlocks or unbounded growth.
		limit = 100
	}
	return &Sink{limit: limit}
}

func (s *Sink) Record(event Event) {
	// Nil receiver check prevents panic if caller circumvents constructor.
	// Empty ID checks act as dead event elimination before acquiring lock.
	if s == nil || event.TraceID == "" || event.EventType == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) == s.limit {
		// O(N) array copy on eviction. Suboptimal for high throughput, but acceptable
		// for small limits. Circular buffer would achieve O(1) amortized.
		copy(s.events, s.events[1:])
		s.events = s.events[:s.limit-1]
	}
	s.events = append(s.events, event)
}
