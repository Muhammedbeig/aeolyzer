package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"path"

	"go.yaml.in/yaml/v3"
)

// EvalSkill is provider-safe skill metadata used to construct a bounded
// candidate set. It intentionally contains no file paths, profiles, action
// classes, tools, connectors, or resource contents.
type EvalSkill struct {
	SkillID           string
	Tier              string
	Description       string
	AntiTriggers      []string
	CompatibleIntents []string
	CapabilityTags    []string
}

// TriggerEvalCase is one validated Layer 4 trigger expectation.
type TriggerEvalCase struct {
	CaseID         string
	TargetSkillID  string
	Group          string
	Input          string
	ExpectedSkill  string
	ForbiddenSkill string
	SafetyCritical bool
}

// TriggerEvalCorpus contains the immutable trigger corpus shipped with the
// process. Checksum identifies the exact registry and fixtures under test.
type TriggerEvalCorpus struct {
	Checksum string
	Skills   []EvalSkill
	Cases    []TriggerEvalCase
}

// LoadEmbeddedTriggerEvalCorpus validates and loads the embedded trigger
// corpus. It fails closed if any case refers to an unregistered skill.
func LoadEmbeddedTriggerEvalCorpus() (TriggerEvalCorpus, error) {
	if _, err := ValidateEmbeddedLibrary(); err != nil {
		return TriggerEvalCorpus{}, err
	}
	return loadTriggerEvalCorpus(embeddedLibrary, ".")
}

func loadTriggerEvalCorpus(fileSystem fs.FS, root string) (TriggerEvalCorpus, error) {
	registryData, err := fs.ReadFile(fileSystem, path.Join(root, "skill-registry.yaml"))
	if err != nil {
		return TriggerEvalCorpus{}, fmt.Errorf("read skill registry: %w", err)
	}
	var registry registryDocument
	if err := yaml.Unmarshal(registryData, &registry); err != nil {
		return TriggerEvalCorpus{}, fmt.Errorf("decode skill registry: %w", err)
	}

	hasher := sha256.New()
	_, _ = hasher.Write(registryData)
	knownSkills := make(map[string]struct{}, len(registry.Skills))
	corpus := TriggerEvalCorpus{
		Skills: make([]EvalSkill, 0, len(registry.Skills)),
	}
	for _, entry := range registry.Skills {
		knownSkills[entry.SkillID] = struct{}{}
		corpus.Skills = append(corpus.Skills, EvalSkill{
			SkillID:           entry.SkillID,
			Tier:              entry.Tier,
			Description:       entry.Description,
			AntiTriggers:      append([]string(nil), entry.AntiTriggers...),
			CompatibleIntents: append([]string(nil), entry.CompatibleIntents...),
			CapabilityTags:    append([]string(nil), entry.CapabilityTags...),
		})

		filename := path.Join(root, entry.Directory, "evals", "trigger_cases.yaml")
		data, err := fs.ReadFile(fileSystem, filename)
		if err != nil {
			return TriggerEvalCorpus{}, fmt.Errorf("read trigger fixture: %w", err)
		}
		_, _ = hasher.Write([]byte{0})
		_, _ = hasher.Write([]byte(entry.SkillID))
		_, _ = hasher.Write([]byte{0})
		_, _ = hasher.Write(data)

		var fixture triggerFixture
		if err := yaml.Unmarshal(data, &fixture); err != nil {
			return TriggerEvalCorpus{}, fmt.Errorf("decode trigger fixture: %w", err)
		}
		appendCases := func(group string, cases []triggerCase, safetyCritical bool) {
			for _, test := range cases {
				expected := test.ExpectedSkill
				if expected == "" && test.ExpectedNotSkill == "" {
					expected = entry.SkillID
				}
				corpus.Cases = append(corpus.Cases, TriggerEvalCase{
					CaseID:         entry.SkillID + "/" + group + "/" + test.ID,
					TargetSkillID:  entry.SkillID,
					Group:          group,
					Input:          test.Input,
					ExpectedSkill:  expected,
					ForbiddenSkill: test.ExpectedNotSkill,
					SafetyCritical: safetyCritical,
				})
			}
		}
		appendCases("positive", fixture.Positive, false)
		appendCases("negative", fixture.Negative, true)
		appendCases("rephrasing", fixture.RephrasingStability, false)
		appendCases("collision", fixture.Collision, false)
		appendCases("out_of_scope", fixture.OutOfScope, true)
	}
	if len(corpus.Skills) == 0 || len(corpus.Cases) == 0 {
		return TriggerEvalCorpus{}, errors.New("embedded trigger eval corpus is empty")
	}
	for _, test := range corpus.Cases {
		if _, found := knownSkills[test.TargetSkillID]; !found {
			return TriggerEvalCorpus{}, errors.New("trigger case target skill is unregistered")
		}
		if test.ExpectedSkill != "" {
			if _, found := knownSkills[test.ExpectedSkill]; !found {
				return TriggerEvalCorpus{}, fmt.Errorf(
					"trigger case %q expects unregistered skill",
					test.CaseID,
				)
			}
		}
		if test.ForbiddenSkill != "" {
			if _, found := knownSkills[test.ForbiddenSkill]; !found {
				return TriggerEvalCorpus{}, fmt.Errorf(
					"trigger case %q forbids unregistered skill",
					test.CaseID,
				)
			}
		}
	}
	corpus.Checksum = "sha256:" + hex.EncodeToString(hasher.Sum(nil))
	return corpus, nil
}
