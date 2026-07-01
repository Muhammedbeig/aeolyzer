package history

import (
	"iter"
	"maps"
	"sync"
	"time"

	"google.golang.org/adk/session"
)

type localSession struct {
	appName   string
	userID    string
	id        string
	mu        sync.RWMutex
	state     map[string]any
	events    []*session.Event
	updatedAt time.Time
}

func (s *localSession) ID() string {
	return s.id
}

func (s *localSession) AppName() string {
	return s.appName
}

func (s *localSession) UserID() string {
	return s.userID
}

func (s *localSession) State() session.State {
	return &localState{session: s}
}

func (s *localSession) Events() session.Events {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copied := make([]*session.Event, len(s.events))
	copy(copied, s.events)
	return localEvents(copied)
}

func (s *localSession) LastUpdateTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}

func (s *localSession) appendEvent(event *session.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	s.updatedAt = event.Timestamp
	for key, value := range event.Actions.StateDelta {
		if isPersistentSessionStateKey(key) {
			s.state[key] = value
		}
	}
}

func (s *localSession) stateSnapshot() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Clone(s.state)
}

type localState struct {
	session *localSession
}

func (s *localState) Get(key string) (any, error) {
	s.session.mu.RLock()
	defer s.session.mu.RUnlock()
	value, ok := s.session.state[key]
	if !ok {
		return nil, session.ErrStateKeyNotExist
	}
	return value, nil
}

func (s *localState) Set(key string, value any) error {
	s.session.mu.Lock()
	defer s.session.mu.Unlock()
	s.session.state[key] = value
	return nil
}

func (s *localState) All() iter.Seq2[string, any] {
	s.session.mu.RLock()
	values := maps.Clone(s.session.state)
	s.session.mu.RUnlock()
	return func(yield func(string, any) bool) {
		for key, value := range values {
			if !yield(key, value) {
				return
			}
		}
	}
}

type localEvents []*session.Event

func (e localEvents) All() iter.Seq[*session.Event] {
	return func(yield func(*session.Event) bool) {
		for _, event := range e {
			if !yield(event) {
				return
			}
		}
	}
}

func (e localEvents) Len() int {
	return len(e)
}

func (e localEvents) At(index int) *session.Event {
	if index < 0 || index >= len(e) {
		return nil
	}
	return e[index]
}
