package mcptransportplane

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("mcp circuit breaker is open")

// CircuitState is the connector/tool circuit state.
type CircuitState string

const (
	// CircuitClosed permits calls.
	CircuitClosed CircuitState = "closed"
	// CircuitOpen denies calls until the reset interval.
	CircuitOpen CircuitState = "open"
	// CircuitHalfOpen permits one probe call.
	CircuitHalfOpen CircuitState = "half_open"
)

// CircuitBreaker is a bounded, concurrency-safe failure breaker.
type CircuitBreaker struct {
	mu               sync.Mutex
	failureThreshold int
	resetAfter       time.Duration
	now              func() time.Time
	state            CircuitState
	failures         int
	openedAt         time.Time
	probeInFlight    bool
}

// NewCircuitBreaker creates a breaker for one connector/tool pair.
func NewCircuitBreaker(
	failureThreshold int,
	resetAfter time.Duration,
	now func() time.Time,
) (*CircuitBreaker, error) {
	if failureThreshold < 1 ||
		failureThreshold > 100 ||
		resetAfter < time.Second ||
		resetAfter > time.Hour ||
		now == nil {
		return nil, errors.New("mcp circuit breaker policy is invalid")
	}
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		resetAfter:       resetAfter,
		now:              now,
		state:            CircuitClosed,
	}, nil
}

// Allow reports whether a call may start.
func (b *CircuitBreaker) Allow() error {
	if b == nil || b.now == nil {
		return ErrCircuitOpen
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == CircuitOpen && b.now().Sub(b.openedAt) >= b.resetAfter {
		b.state = CircuitHalfOpen
		b.probeInFlight = false
	}
	switch b.state {
	case CircuitClosed:
		return nil
	case CircuitHalfOpen:
		if b.probeInFlight {
			return ErrCircuitOpen
		}
		b.probeInFlight = true
		return nil
	default:
		return ErrCircuitOpen
	}
}

// Success records a successful call and closes the circuit.
func (b *CircuitBreaker) Success() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = CircuitClosed
	b.failures = 0
	b.probeInFlight = false
}

// Failure records a failed call and opens the circuit at threshold.
func (b *CircuitBreaker) Failure() {
	if b == nil || b.now == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.probeInFlight = false
	if b.state == CircuitHalfOpen {
		b.state = CircuitOpen
		b.openedAt = b.now()
		return
	}
	b.failures++
	if b.failures >= b.failureThreshold {
		b.state = CircuitOpen
		b.openedAt = b.now()
	}
}

// State returns the current state without changing it.
func (b *CircuitBreaker) State() CircuitState {
	if b == nil {
		return CircuitOpen
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
