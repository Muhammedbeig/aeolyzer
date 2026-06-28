package observability

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSinkRemainsBoundedUnderConcurrentWrites(t *testing.T) {
	t.Parallel()

	const limit = 20
	sink := NewSink(limit)
	var waitGroup sync.WaitGroup
	for index := 0; index < 100; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			sink.Record(Event{
				TraceID:   fmt.Sprintf("trace-%d", index),
				EventType: "test",
				Outcome:   "succeeded",
				At:        time.Unix(int64(index), 0),
			})
		}(index)
	}
	waitGroup.Wait()

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.events) != limit {
		t.Fatalf("events = %d, want %d", len(sink.events), limit)
	}
}
