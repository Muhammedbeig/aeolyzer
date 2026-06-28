package skills

import (
	"strings"
	"testing"
)

func TestLoadEmbeddedTriggerEvalCorpus(t *testing.T) {
	corpus, err := LoadEmbeddedTriggerEvalCorpus()
	if err != nil {
		t.Fatalf("LoadEmbeddedTriggerEvalCorpus() error = %v", err)
	}
	if len(corpus.Skills) != 44 {
		t.Fatalf("skills = %d, want 44", len(corpus.Skills))
	}
	if len(corpus.Cases) != 484 {
		t.Fatalf("cases = %d, want 484", len(corpus.Cases))
	}
	if !strings.HasPrefix(corpus.Checksum, "sha256:") {
		t.Fatalf("checksum = %q, want sha256 prefix", corpus.Checksum)
	}
	for _, skill := range corpus.Skills {
		if skill.Description == "" || len(skill.CompatibleIntents) == 0 {
			t.Fatalf("skill metadata is incomplete: %+v", skill)
		}
	}
}
