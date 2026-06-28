package agbom_test

import (
	"aeolyzer/internal/observability/agbom"
	"testing"
)

func TestBuildAgBOM(t *testing.T) {
	_, err := agbom.BuildAgBOM("")
	if err == nil {
		t.Fatal("expected trace ID requirement for agbom")
	}
}
