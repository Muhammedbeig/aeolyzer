package feedbackimprovementloop

import "testing"

func TestMineCorrectionsRequiresRepeatedCrossTenantEvidence(t *testing.T) {
	corrections := []Correction{
		{
			FailureClass:    "unsupported_claim",
			Component:       "content_research",
			CorrectionClass: "add_source_requirement",
			TenantHash:      "tenant-a",
		},
		{
			FailureClass:    "unsupported_claim",
			Component:       "content_research",
			CorrectionClass: "add_source_requirement",
			TenantHash:      "tenant-b",
		},
		{
			FailureClass:    "unsupported_claim",
			Component:       "content_research",
			CorrectionClass: "add_source_requirement",
			TenantHash:      "tenant-c",
		},
	}
	recommendations, err := MineCorrections(corrections, 3, 3)
	if err != nil {
		t.Fatalf("MineCorrections() failed: %v", err)
	}
	if len(recommendations) != 1 ||
		!recommendations[0].RequiresHumanReview {
		t.Fatalf("MineCorrections() = %+v, want one reviewed recommendation", recommendations)
	}
	recommendations, err = MineCorrections(corrections[:2], 3, 2)
	if err != nil {
		t.Fatalf("MineCorrections(insufficient) failed: %v", err)
	}
	if len(recommendations) != 0 {
		t.Fatalf("MineCorrections(insufficient) = %+v, want none", recommendations)
	}
}
