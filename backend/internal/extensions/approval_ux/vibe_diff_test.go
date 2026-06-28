package approval_ux_test

import (
	"testing"
	"aeolyzer/internal/extensions/approval_ux"
)

func TestVibeDiffMetadataLeak(t *testing.T) {
	diff := approval_ux.VibeDiff{
		Summary: "Update constraints",
		RiskNotes: []string{"trace_id"},
	}
	if err := approval_ux.ValidateVibeDiff(diff); err == nil {
		t.Fatal("expected internal metadata leak to be blocked")
	}
}
