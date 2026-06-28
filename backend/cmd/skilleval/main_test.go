package main

import (
	"testing"

	observabilityconfig "aeolyzer/internal/observability/config"
	"aeolyzer/internal/skills"
)

func TestBuildSkillEvalBatchesCoversCorpus(t *testing.T) {
	corpus, err := skills.LoadEmbeddedTriggerEvalCorpus()
	if err != nil {
		t.Fatalf("LoadEmbeddedTriggerEvalCorpus() error = %v", err)
	}
	batches, err := buildSkillEvalBatches(
		corpus,
		"test-model",
		"test-prompt",
		"",
	)
	if err != nil {
		t.Fatalf("buildSkillEvalBatches() error = %v", err)
	}
	if len(batches) != 44 {
		t.Fatalf("batches = %d, want 44", len(batches))
	}
	totalCases := 0
	for _, batch := range batches {
		if len(batch.batch.Candidates) != defaultCandidateCount {
			t.Fatalf(
				"skill %s candidates = %d, want %d",
				batch.batch.TargetSkillID,
				len(batch.batch.Candidates),
				defaultCandidateCount,
			)
		}
		if len(batch.batch.Cases) != 11 {
			t.Fatalf(
				"skill %s cases = %d, want 11",
				batch.batch.TargetSkillID,
				len(batch.batch.Cases),
			)
		}
		totalCases += len(batch.batch.Cases)
	}
	if totalCases != 484 {
		t.Fatalf("total cases = %d, want 484", totalCases)
	}
}

func TestPassKForTierUsesProductionPolicy(t *testing.T) {
	policies, err := observabilityconfig.LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() error = %v", err)
	}
	tests := map[string]int{
		"read":  policies.Eval.PassK.ReadWorkflow,
		"draft": policies.Eval.PassK.DraftWorkflow,
		"act":   policies.Eval.PassK.GuardedWriteWorkflow,
	}
	for tier, expected := range tests {
		actual, err := passKForTier(tier, policies.Eval)
		if err != nil {
			t.Fatalf("passKForTier(%q) error = %v", tier, err)
		}
		if actual != expected {
			t.Fatalf("passKForTier(%q) = %d, want %d", tier, actual, expected)
		}
	}
}
