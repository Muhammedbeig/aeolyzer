package skills

import (
	"path/filepath"
	"testing"
)

func TestRepositorySkillLibraryValidates(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	report, err := ValidateLibrary(root)
	if err != nil {
		t.Fatalf("ValidateLibrary() failed: %v", err)
	}
	if report.SkillsValidated != 44 {
		t.Fatalf("skills validated = %d, want 44", report.SkillsValidated)
	}
	if report.PositiveTriggerCases < 132 ||
		report.NegativeTriggerCases < 132 ||
		report.RephrasingCases < 132 ||
		report.GoldenCases < 44 ||
		report.TrajectoryCases < 44 ||
		report.RegressionCases < 132 {
		t.Fatalf("eval corpus is below minimum: %+v", report)
	}
}

func TestEmbeddedSkillLibraryValidates(t *testing.T) {
	report, err := ValidateEmbeddedLibrary()
	if err != nil {
		t.Fatalf("ValidateEmbeddedLibrary() failed: %v", err)
	}
	if report.SkillsValidated != 44 {
		t.Fatalf("skills validated = %d, want 44", report.SkillsValidated)
	}
}
