package releasegate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFindsMissingProductionEvidence(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "internal/skills/skill-registry.yaml", "version: 2\nskills: []\n")
	writeTestFile(t, root, "internal/skills/skill.schema.json", "{}")
	writeTestFile(t, root, "internal/skills/skills/example/SKILL.md", `---
name: example
description: Example skill. Do NOT use outside tests.
---

# Example
`)

	report, err := Check(root)
	if err != nil {
		t.Fatalf("Check(%q) returned unexpected error: %v", root, err)
	}
	if report.Ready() {
		t.Fatalf("Check(%q).Ready() = true, want false", root)
	}

	codes := make(map[string]bool)
	for _, finding := range report.Findings {
		codes[finding.Code] = true
	}
	for _, code := range []string{
		"placeholder_artifact",
		"skill_registry_empty",
		"skill_frontmatter_incomplete",
		"skill_sections_incomplete",
		"skill_eval_evidence_missing",
		"executable_control_missing",
	} {
		if !codes[code] {
			t.Errorf("Check(%q) finding codes = %v, want %q", root, codes, code)
		}
	}
}

func TestCheckRejectsFileAsRoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "repository")
	if err := os.WriteFile(path, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) failed: %v", path, err)
	}

	if _, err := Check(path); err == nil {
		t.Errorf("Check(%q) returned nil error, want non-nil", path)
	}
}

func TestParseSkillFile(t *testing.T) {
	fields, body, ok := parseSkillFile(`---
name: test-skill
description: |
  Tests parsing.
version: 1.0.0
---

# Test Skill

## Purpose
Exercise the parser.
`)
	if !ok {
		t.Fatal("parseSkillFile() = ok false, want true")
	}
	for _, field := range []string{"name", "description", "version"} {
		if _, exists := fields[field]; !exists {
			t.Errorf("parseSkillFile() fields missing %q", field)
		}
	}
	if !hasMarkdownSection(body, "Purpose") {
		t.Error("hasMarkdownSection(body, \"Purpose\") = false, want true")
	}
}

func TestRegistrySkillStatuses(t *testing.T) {
	statuses := registrySkillStatuses(`skills:
  - skill_id: topic_discovery
    name: topic-discovery
    status: experimental
  - skill_id: writing
    name: writing
    status: active
`)
	if statuses["topic_discovery"] != "experimental" {
		t.Fatalf(
			"registrySkillStatuses(topic_discovery) = %q",
			statuses["topic_discovery"],
		)
	}
	if statuses["writing"] != "active" {
		t.Fatalf("registrySkillStatuses(writing) = %q", statuses["writing"])
	}
}

func writeTestFile(t *testing.T, root, relativePath, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) failed: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) failed: %v", path, err)
	}
}
