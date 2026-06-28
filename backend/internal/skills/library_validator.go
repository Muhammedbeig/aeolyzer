package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"strings"

	"go.yaml.in/yaml/v3"
)

// LibraryReport contains deterministic Layer 4 validation evidence.
type LibraryReport struct {
	SkillsValidated      int `json:"skills_validated"`
	PositiveTriggerCases int `json:"positive_trigger_cases"`
	NegativeTriggerCases int `json:"negative_trigger_cases"`
	RephrasingCases      int `json:"rephrasing_cases"`
	GoldenCases          int `json:"golden_cases"`
	TrajectoryCases      int `json:"trajectory_cases"`
	RegressionCases      int `json:"regression_cases"`
}

type registryDocument struct {
	Skills []registryEntry `yaml:"skills"`
}

type registryEntry struct {
	SkillID           string   `yaml:"skill_id"`
	Directory         string   `yaml:"directory"`
	Status            string   `yaml:"status"`
	Tier              string   `yaml:"tier"`
	Description       string   `yaml:"description"`
	AntiTriggers      []string `yaml:"anti_triggers"`
	CompatibleIntents []string `yaml:"compatible_intents"`
	CapabilityTags    []string `yaml:"capability_tags"`
	Checksum          string   `yaml:"checksum"`
}

type triggerFixture struct {
	Version             int           `yaml:"version"`
	SkillID             string        `yaml:"skill_id"`
	Positive            []triggerCase `yaml:"positive"`
	Negative            []triggerCase `yaml:"negative"`
	RephrasingStability []triggerCase `yaml:"rephrasing_stability"`
	Collision           []triggerCase `yaml:"collision"`
	OutOfScope          []triggerCase `yaml:"out_of_scope"`
}

type triggerCase struct {
	ID               string `yaml:"id"`
	SourceCase       string `yaml:"source_case"`
	Input            string `yaml:"input"`
	ExpectedSkill    string `yaml:"expected_skill"`
	ExpectedNotSkill string `yaml:"expected_not_skill"`
	ExpectedIntent   string `yaml:"expected_intent"`
}

type goldenFixture struct {
	Version int          `yaml:"version"`
	SkillID string       `yaml:"skill_id"`
	Cases   []goldenCase `yaml:"cases"`
}

type goldenCase struct {
	ID                 string   `yaml:"id"`
	Input              string   `yaml:"input"`
	RequiredOutputs    []string `yaml:"required_outputs"`
	RequiredQualities  []string `yaml:"required_qualities"`
	ForbiddenQualities []string `yaml:"forbidden_qualities"`
}

type trajectoryFixture struct {
	Version int              `yaml:"version"`
	SkillID string           `yaml:"skill_id"`
	Cases   []trajectoryCase `yaml:"cases"`
}

type trajectoryCase struct {
	ID                     string   `yaml:"id"`
	Mode                   string   `yaml:"mode"`
	ExpectedActionClasses  []string `yaml:"expected_action_classes"`
	ForbiddenActionClasses []string `yaml:"forbidden_action_classes"`
}

type rubricFixture struct {
	Version    int               `yaml:"version"`
	RubricID   string            `yaml:"rubric_id"`
	SkillID    string            `yaml:"skill_id"`
	PassScore  float64           `yaml:"pass_score"`
	Dimensions []rubricDimension `yaml:"dimensions"`
}

type rubricDimension struct {
	ID      string  `yaml:"id"`
	Weight  float64 `yaml:"weight"`
	Minimum int     `yaml:"minimum"`
	Maximum int     `yaml:"maximum"`
}

type regressionFixture struct {
	Version int              `yaml:"version"`
	SkillID string           `yaml:"skill_id"`
	Cases   []regressionCase `yaml:"cases"`
}

type regressionCase struct {
	ID       string `yaml:"id"`
	Input    string `yaml:"input"`
	Expected string `yaml:"expected"`
}

// ValidateLibrary validates the registry, checksums, skill contracts,
// manifests, and all static evaluation fixtures.
func ValidateLibrary(root string) (LibraryReport, error) {
	return validateLibraryFS(os.DirFS(root), "internal/skills")
}

// ValidateEmbeddedLibrary validates the exact skill artifacts shipped in the
// compiled binary. Startup can therefore fail closed without relying on the
// process working directory or mutable deployment files.
func ValidateEmbeddedLibrary() (LibraryReport, error) {
	return validateLibraryFS(embeddedLibrary, ".")
}

func validateLibraryFS(fileSystem fs.FS, root string) (LibraryReport, error) {
	schemas, err := NewSchemas()
	if err != nil {
		return LibraryReport{}, err
	}
	registryData, err := fs.ReadFile(fileSystem, path.Join(root, "skill-registry.yaml"))
	if err != nil {
		return LibraryReport{}, fmt.Errorf("read skill registry: %w", err)
	}
	if err := schemas.ValidateRegistry(registryData); err != nil {
		return LibraryReport{}, fmt.Errorf("validate skill registry schema: %w", err)
	}
	var registry registryDocument
	if err := yaml.Unmarshal(registryData, &registry); err != nil {
		return LibraryReport{}, fmt.Errorf("decode skill registry: %w", err)
	}
	if len(registry.Skills) == 0 {
		return LibraryReport{}, errors.New("skill registry is empty")
	}

	seen := make(map[string]struct{}, len(registry.Skills))
	var report LibraryReport
	for _, entry := range registry.Skills {
		if _, duplicate := seen[entry.SkillID]; duplicate {
			return LibraryReport{}, fmt.Errorf("duplicate skill %q", entry.SkillID)
		}
		seen[entry.SkillID] = struct{}{}
		if entry.Status == "active" {
			return LibraryReport{}, fmt.Errorf(
				"skill %q is active without recorded Layer 8 release evidence",
				entry.SkillID,
			)
		}
		if err := validateSkillDirectory(fileSystem, root, schemas, entry, &report); err != nil {
			return LibraryReport{}, fmt.Errorf("validate skill %s: %w", entry.SkillID, err)
		}
		report.SkillsValidated++
	}
	return report, nil
}

func validateSkillDirectory(
	fileSystem fs.FS,
	root string,
	schemas *Schemas,
	entry registryEntry,
	report *LibraryReport,
) error {
	directory := path.Join(root, entry.Directory)
	skillData, err := fs.ReadFile(fileSystem, path.Join(directory, "SKILL.md"))
	if err != nil {
		return err
	}
	digest := sha256.Sum256(skillData)
	actualChecksum := "sha256:" + hex.EncodeToString(digest[:])
	if actualChecksum != entry.Checksum {
		return ErrChecksumMismatch
	}
	frontmatter, body, err := splitSkillFile(skillData)
	if err != nil {
		return err
	}
	if err := schemas.ValidateSkillFrontmatter(frontmatter); err != nil {
		return err
	}
	for _, section := range []string{
		"Purpose",
		"When to use",
		"When NOT to use",
		"Inputs expected",
		"Procedure",
		"Output contract",
		"Quality gates",
		"Boundary rules",
		"Resources",
		"Failure behavior",
	} {
		if !containsHeading(body, section) {
			return fmt.Errorf("missing section %q", section)
		}
	}
	if len(strings.Fields(body)) > 5000 {
		return errors.New("skill body exceeds hard word limit")
	}

	resourceData, err := fs.ReadFile(fileSystem, path.Join(directory, "resource-manifest.yaml"))
	if err != nil {
		return err
	}
	if err := schemas.ValidateResourceManifest(resourceData); err != nil {
		return err
	}
	evalManifestData, err := fs.ReadFile(fileSystem, path.Join(directory, "eval-manifest.yaml"))
	if err != nil {
		return err
	}
	if err := schemas.ValidateEvalManifest(evalManifestData); err != nil {
		return err
	}
	for _, required := range []string{"OWNERS", "CHANGELOG.md"} {
		info, err := fs.Stat(fileSystem, path.Join(directory, required))
		if err != nil || info.Size() == 0 {
			return fmt.Errorf("%s is missing or empty", required)
		}
	}
	return validateEvalFixtures(fileSystem, directory, entry, report)
}

func validateEvalFixtures(
	fileSystem fs.FS,
	directory string,
	entry registryEntry,
	report *LibraryReport,
) error {
	var triggers triggerFixture
	if err := decodeYAMLFile(fileSystem, path.Join(directory, "evals", "trigger_cases.yaml"), &triggers); err != nil {
		return err
	}
	if triggers.Version != 1 ||
		triggers.SkillID != entry.SkillID ||
		len(triggers.Positive) < 3 ||
		len(triggers.Negative) < 3 ||
		len(triggers.RephrasingStability) < len(triggers.Positive) ||
		len(triggers.Collision) < 1 ||
		len(triggers.OutOfScope) < 1 {
		return errors.New("trigger fixture does not meet minimum corpus")
	}
	if err := validateTriggerCases(entry.SkillID, triggers); err != nil {
		return err
	}

	var golden goldenFixture
	if err := decodeYAMLFile(fileSystem, path.Join(directory, "evals", "golden_cases.yaml"), &golden); err != nil {
		return err
	}
	if golden.Version != 1 || golden.SkillID != entry.SkillID || len(golden.Cases) < 1 {
		return errors.New("golden fixture is incomplete")
	}
	for _, test := range golden.Cases {
		if test.ID == "" ||
			test.Input == "" ||
			len(test.RequiredOutputs) == 0 ||
			len(test.RequiredQualities) < 3 ||
			len(test.ForbiddenQualities) < 3 {
			return errors.New("golden case is incomplete")
		}
	}

	var trajectories trajectoryFixture
	if err := decodeYAMLFile(fileSystem, path.Join(directory, "evals", "trajectory_cases.yaml"), &trajectories); err != nil {
		return err
	}
	if trajectories.Version != 1 ||
		trajectories.SkillID != entry.SkillID ||
		len(trajectories.Cases) < 1 {
		return errors.New("trajectory fixture is incomplete")
	}
	expectedMode := "ANY_ORDER"
	if entry.Tier == "draft" {
		expectedMode = "IN_ORDER"
	} else if entry.Tier == "act" {
		expectedMode = "EXACT"
	}
	for _, test := range trajectories.Cases {
		if test.ID == "" ||
			test.Mode != expectedMode ||
			len(test.ExpectedActionClasses) == 0 ||
			len(test.ForbiddenActionClasses) < 3 {
			return errors.New("trajectory case violates tier policy")
		}
	}

	var rubric rubricFixture
	if err := decodeYAMLFile(fileSystem, path.Join(directory, "evals", "rubric.yaml"), &rubric); err != nil {
		return err
	}
	if rubric.Version != 1 ||
		rubric.SkillID != entry.SkillID ||
		rubric.RubricID == "" ||
		rubric.PassScore < 4 ||
		len(rubric.Dimensions) < 5 {
		return errors.New("rubric fixture is incomplete")
	}
	weight := 0.0
	for _, dimension := range rubric.Dimensions {
		if dimension.ID == "" ||
			dimension.Weight <= 0 ||
			dimension.Minimum != 1 ||
			dimension.Maximum != 5 {
			return errors.New("rubric dimension is invalid")
		}
		weight += dimension.Weight
	}
	if math.Abs(weight-1) > 0.000001 {
		return errors.New("rubric weights must total one")
	}

	var regressions regressionFixture
	if err := decodeYAMLFile(fileSystem, path.Join(directory, "evals", "regression_cases.yaml"), &regressions); err != nil {
		return err
	}
	if regressions.Version != 1 ||
		regressions.SkillID != entry.SkillID ||
		len(regressions.Cases) < 3 {
		return errors.New("regression fixture is incomplete")
	}
	for _, test := range regressions.Cases {
		if test.ID == "" || test.Input == "" || test.Expected == "" {
			return errors.New("regression case is incomplete")
		}
	}

	report.PositiveTriggerCases += len(triggers.Positive)
	report.NegativeTriggerCases += len(triggers.Negative)
	report.RephrasingCases += len(triggers.RephrasingStability)
	report.GoldenCases += len(golden.Cases)
	report.TrajectoryCases += len(trajectories.Cases)
	report.RegressionCases += len(regressions.Cases)
	return nil
}

func validateTriggerCases(skillID string, fixture triggerFixture) error {
	seenIDs := make(map[string]struct{})
	for _, test := range fixture.Positive {
		if test.ID == "" ||
			test.Input == "" ||
			test.ExpectedSkill != skillID ||
			test.ExpectedIntent == "" {
			return errors.New("positive trigger case is invalid")
		}
		if _, duplicate := seenIDs[test.ID]; duplicate {
			return errors.New("trigger case id is duplicated")
		}
		seenIDs[test.ID] = struct{}{}
	}
	for _, test := range fixture.Negative {
		if test.ID == "" ||
			test.Input == "" ||
			test.ExpectedNotSkill != skillID {
			return errors.New("negative trigger case is invalid")
		}
		if _, duplicate := seenIDs[test.ID]; duplicate {
			return errors.New("trigger case id is duplicated")
		}
		seenIDs[test.ID] = struct{}{}
	}
	for _, test := range append(
		append([]triggerCase(nil), fixture.Collision...),
		fixture.OutOfScope...,
	) {
		if test.ID == "" || test.Input == "" || test.ExpectedNotSkill != skillID {
			return errors.New("collision or out-of-scope case is invalid")
		}
	}
	for _, test := range fixture.RephrasingStability {
		if test.SourceCase == "" ||
			test.Input == "" ||
			test.ExpectedSkill != skillID {
			return errors.New("rephrasing stability case is invalid")
		}
	}
	return nil
}

func splitSkillFile(data []byte) ([]byte, string, error) {
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return nil, "", errors.New("skill frontmatter is missing")
	}
	parts := strings.SplitN(strings.TrimPrefix(normalized, "---\n"), "\n---\n", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("skill frontmatter is unterminated")
	}
	return []byte(parts[0]), parts[1], nil
}

func containsHeading(body, heading string) bool {
	for _, line := range strings.Split(body, "\n") {
		if strings.EqualFold(strings.TrimSpace(line), "## "+heading) {
			return true
		}
	}
	return false
}

func decodeYAMLFile(fileSystem fs.FS, filename string, destination any) error {
	data, err := fs.ReadFile(fileSystem, filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("decode %s: %w", path.Base(filename), err)
	}
	return nil
}
