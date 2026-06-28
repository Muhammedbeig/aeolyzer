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
	mu     sync.Mutex
	events []Event
	limit  int
}

func NewSink(limit int) *Sink {
	if limit < 1 {
		limit = 100
	}
	return &Sink{limit: limit}
}

func (s *Sink) Record(event Event) {
	if s == nil || event.TraceID == "" || event.EventType == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) == s.limit {
		copy(s.events, s.events[1:])
		s.events = s.events[:s.limit-1]
	}
	s.events = append(s.events, event)
}
