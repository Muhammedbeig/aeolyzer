package observability

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSinkRemainsBoundedUnderConcurrentWrites(t *testing.T) {
	// Ensures lock contention and data races are reliably exposed during test execution.
	t.Parallel()

	const limit = 20
	sink := NewSink(limit)
	var waitGroup sync.WaitGroup
	// 100 concurrent producers against a limit of 20 forces aggressive eviction testing.
	for index := 0; index < 100; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			sink.Record(Event{
				TraceID:   fmt.Sprintf("trace-%d", index),
				EventType: "test",
				Outcome:   "succeeded",
				// Deterministic timestamp injection to avoid flakey state behavior.
				At:        time.Unix(int64(index), 0),
			})
		}(index)
	}
	waitGroup.Wait()

	// Direct lock inspection bypasses standard APIs to verify internal structure invariants.
	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.events) != limit {
		t.Fatalf("events = %d, want %d", len(sink.events), limit)
	}
}
