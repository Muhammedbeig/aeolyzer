package evaluationengine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

const (
	minSkillEvalCandidates = 5
	maxSkillEvalCandidates = 15
	maxSkillEvalCases      = 32
	maxSkillEvalPassK      = 10
)

// SkillEvalCandidate is one bounded candidate. SkillID remains local; the
// provider receives only a run-specific opaque alias.
type SkillEvalCandidate struct {
	SkillID      string
	Description  string
	AntiTriggers []string
}

// SkillEvalCase is one expected routing decision.
type SkillEvalCase struct {
	CaseID         string
	Group          string
	Input          string
	ExpectedSkill  string
	ForbiddenSkill string
	SafetyCritical bool
}

// SkillEvalBatch is one target skill's bounded evaluation batch.
type SkillEvalBatch struct {
	EvalID         string
	CorpusChecksum string
	TargetSkillID  string
	Model          string
	PromptVersion  string
	Candidates     []SkillEvalCandidate
	Cases          []SkillEvalCase
}

// SkillSelectionCandidate is provider-safe candidate metadata.
type SkillSelectionCandidate struct {
	Alias        string   `json:"alias"`
	Description  string   `json:"description"`
	AntiTriggers []string `json:"anti_triggers"`
}

// SkillSelectionCase is provider-safe eval input without expected answers.
type SkillSelectionCase struct {
	CaseID string `json:"case_id"`
	Input  string `json:"input"`
}

// SkillSelectionRequest is the provider-neutral structured selection request.
type SkillSelectionRequest struct {
	Model         string                    `json:"model"`
	PromptVersion string                    `json:"prompt_version"`
	EvalID        string                    `json:"eval_id"`
	Candidates    []SkillSelectionCandidate `json:"candidates"`
	Cases         []SkillSelectionCase      `json:"cases"`
}

// SkillSelectionClient calls a Layer 7 model provider adapter.
type SkillSelectionClient interface {
	SelectSkills(context.Context, SkillSelectionRequest) ([]byte, error)
}

// SkillEvalCaseResult contains sanitized per-case evidence.
type SkillEvalCaseResult struct {
	CaseHash          string  `json:"case_hash"`
	Group             string  `json:"group"`
	Passed            bool    `json:"passed"`
	Flaky             bool    `json:"flaky"`
	MinimumConfidence float64 `json:"minimum_confidence"`
}

// SkillEvalReport is a sanitized, deterministic release-evidence record. It
// intentionally excludes raw test prompts and provider summaries.
type SkillEvalReport struct {
	EvalID            string                `json:"eval_id"`
	CorpusChecksum    string                `json:"corpus_checksum"`
	TargetSkillID     string                `json:"target_skill_id"`
	Model             string                `json:"model"`
	PromptVersion     string                `json:"prompt_version"`
	Runs              int                   `json:"runs"`
	TotalCases        int                   `json:"total_cases"`
	PassedCases       int                   `json:"passed_cases"`
	FailedCases       int                   `json:"failed_cases"`
	SafetyFailures    int                   `json:"safety_failures"`
	FlakyCases        int                   `json:"flaky_cases"`
	FlakeRate         float64               `json:"flake_rate"`
	MinimumConfidence float64               `json:"minimum_confidence"`
	Passed            bool                  `json:"passed"`
	Cases             []SkillEvalCaseResult `json:"cases"`
}

type skillSelectionWireOutput struct {
	Results []struct {
		CaseID            string  `json:"case_id"`
		SelectedCandidate string  `json:"selected_candidate"`
		Confidence        float64 `json:"confidence"`
		Summary           string  `json:"summary"`
	} `json:"results"`
}

type skillEvalAccumulator struct {
	passed        bool
	minConfidence float64
	selections    []string
}

// EvaluateSkillTriggers runs pass^k skill routing evaluation. Candidate order
// and aliases are rotated on every run to expose position-sensitive routing.
// Every run must pass; any safety-case failure blocks the batch.
func EvaluateSkillTriggers(
	ctx context.Context,
	client SkillSelectionClient,
	batch SkillEvalBatch,
	passK int,
	minimumConfidence float64,
) (SkillEvalReport, error) {
	if client == nil {
		return SkillEvalReport{}, errors.New("skill selection client is required")
	}
	if err := validateSkillEvalBatch(batch, passK, minimumConfidence); err != nil {
		return SkillEvalReport{}, err
	}

	accumulators := make(map[string]*skillEvalAccumulator, len(batch.Cases))
	for _, test := range batch.Cases {
		accumulators[test.CaseID] = &skillEvalAccumulator{
			passed:        true,
			minConfidence: 1,
			selections:    make([]string, 0, passK),
		}
	}

	for run := 0; run < passK; run++ {
		if err := ctx.Err(); err != nil {
			return SkillEvalReport{}, fmt.Errorf("evaluate skill triggers: %w", err)
		}
		request, aliasToSkill := buildSkillSelectionRequest(batch, run)
		raw, err := client.SelectSkills(ctx, request)
		if err != nil {
			return SkillEvalReport{}, fmt.Errorf("select skills on run %d: %w", run+1, err)
		}
		results, err := validateSkillSelectionOutput(raw, request)
		if err != nil {
			return SkillEvalReport{}, fmt.Errorf(
				"validate skill selection output on run %d: %w",
				run+1,
				err,
			)
		}
		for _, test := range batch.Cases {
			result := results[test.CaseID]
			selectedSkill := ""
			if result.SelectedCandidate != "none" {
				selectedSkill = aliasToSkill[result.SelectedCandidate]
			}
			correct := selectedSkill == test.ExpectedSkill
			if test.ForbiddenSkill != "" && selectedSkill == test.ForbiddenSkill {
				correct = false
			}
			if result.Confidence < minimumConfidence {
				correct = false
			}
			accumulator := accumulators[test.CaseID]
			accumulator.passed = accumulator.passed && correct
			accumulator.minConfidence = math.Min(
				accumulator.minConfidence,
				result.Confidence,
			)
			accumulator.selections = append(accumulator.selections, selectedSkill)
		}
	}

	report := SkillEvalReport{
		EvalID:            batch.EvalID,
		CorpusChecksum:    batch.CorpusChecksum,
		TargetSkillID:     batch.TargetSkillID,
		Model:             batch.Model,
		PromptVersion:     batch.PromptVersion,
		Runs:              passK,
		TotalCases:        len(batch.Cases),
		MinimumConfidence: minimumConfidence,
		Passed:            true,
		Cases:             make([]SkillEvalCaseResult, 0, len(batch.Cases)),
	}
	for _, test := range batch.Cases {
		accumulator := accumulators[test.CaseID]
		flaky := hasDifferentStrings(accumulator.selections)
		report.Cases = append(report.Cases, SkillEvalCaseResult{
			CaseHash:          hashSkillEvalCase(test),
			Group:             test.Group,
			Passed:            accumulator.passed,
			Flaky:             flaky,
			MinimumConfidence: accumulator.minConfidence,
		})
		if accumulator.passed {
			report.PassedCases++
		} else {
			report.FailedCases++
			report.Passed = false
			if test.SafetyCritical {
				report.SafetyFailures++
			}
		}
		if flaky {
			report.FlakyCases++
			report.Passed = false
		}
	}
	report.FlakeRate = float64(report.FlakyCases) / float64(report.TotalCases)
	return report, nil
}

// ValidateSkillSelectionRequest validates the Layer 8 request contract before
// a Layer 7 provider adapter transmits it.
func ValidateSkillSelectionRequest(request SkillSelectionRequest) error {
	if request.Model == "" || request.PromptVersion == "" || request.EvalID == "" {
		return errors.New("skill selection model, prompt version, and eval id are required")
	}
	if len(request.Candidates) < minSkillEvalCandidates ||
		len(request.Candidates) > maxSkillEvalCandidates {
		return errors.New("skill selection candidate count is invalid")
	}
	if len(request.Cases) == 0 || len(request.Cases) > maxSkillEvalCases {
		return errors.New("skill selection case count is invalid")
	}
	aliases := make(map[string]struct{}, len(request.Candidates))
	for _, candidate := range request.Candidates {
		if candidate.Alias == "" ||
			len(candidate.Description) < 20 ||
			len(candidate.Description) > 720 ||
			len(candidate.AntiTriggers) == 0 ||
			len(candidate.AntiTriggers) > 16 {
			return errors.New("skill selection candidate is invalid")
		}
		if _, duplicate := aliases[candidate.Alias]; duplicate {
			return errors.New("skill selection candidate alias is duplicated")
		}
		aliases[candidate.Alias] = struct{}{}
	}
	caseIDs := make(map[string]struct{}, len(request.Cases))
	for _, test := range request.Cases {
		if test.CaseID == "" || test.Input == "" || len(test.Input) > 4096 {
			return errors.New("skill selection case is invalid")
		}
		if _, duplicate := caseIDs[test.CaseID]; duplicate {
			return errors.New("skill selection case id is duplicated")
		}
		caseIDs[test.CaseID] = struct{}{}
	}
	return nil
}

func validateSkillEvalBatch(
	batch SkillEvalBatch,
	passK int,
	minimumConfidence float64,
) error {
	if batch.EvalID == "" ||
		batch.CorpusChecksum == "" ||
		batch.TargetSkillID == "" ||
		batch.Model == "" ||
		batch.PromptVersion == "" {
		return errors.New("skill eval batch metadata is incomplete")
	}
	if passK < 1 || passK > maxSkillEvalPassK {
		return errors.New("skill eval pass k is invalid")
	}
	if minimumConfidence < 0.5 || minimumConfidence > 1 {
		return errors.New("skill eval minimum confidence is invalid")
	}
	if len(batch.Candidates) < minSkillEvalCandidates ||
		len(batch.Candidates) > maxSkillEvalCandidates ||
		len(batch.Cases) == 0 ||
		len(batch.Cases) > maxSkillEvalCases {
		return errors.New("skill eval batch size is invalid")
	}
	skills := make(map[string]struct{}, len(batch.Candidates))
	for _, candidate := range batch.Candidates {
		if candidate.SkillID == "" ||
			len(candidate.Description) < 20 ||
			len(candidate.AntiTriggers) == 0 {
			return errors.New("skill eval candidate is incomplete")
		}
		if _, duplicate := skills[candidate.SkillID]; duplicate {
			return errors.New("skill eval candidate is duplicated")
		}
		skills[candidate.SkillID] = struct{}{}
	}
	if _, found := skills[batch.TargetSkillID]; !found {
		return errors.New("skill eval target is absent from candidate set")
	}
	caseIDs := make(map[string]struct{}, len(batch.Cases))
	for _, test := range batch.Cases {
		if test.CaseID == "" ||
			test.Group == "" ||
			test.Input == "" ||
			len(test.Input) > 4096 {
			return errors.New("skill eval case is incomplete")
		}
		if _, duplicate := caseIDs[test.CaseID]; duplicate {
			return errors.New("skill eval case id is duplicated")
		}
		caseIDs[test.CaseID] = struct{}{}
		if test.ExpectedSkill != "" {
			if _, found := skills[test.ExpectedSkill]; !found {
				return errors.New("expected skill is absent from candidate set")
			}
		}
		if test.ForbiddenSkill != "" {
			if _, found := skills[test.ForbiddenSkill]; !found {
				return errors.New("forbidden skill is absent from candidate set")
			}
		}
	}
	return nil
}

func buildSkillSelectionRequest(
	batch SkillEvalBatch,
	run int,
) (SkillSelectionRequest, map[string]string) {
	count := len(batch.Candidates)
	candidates := make([]SkillSelectionCandidate, 0, count)
	aliasToSkill := make(map[string]string, count)
	for index := 0; index < count; index++ {
		candidate := batch.Candidates[(index+run)%count]
		alias := "candidate_" + strconv.Itoa(index+1)
		candidates = append(candidates, SkillSelectionCandidate{
			Alias:        alias,
			Description:  candidate.Description,
			AntiTriggers: append([]string(nil), candidate.AntiTriggers...),
		})
		aliasToSkill[alias] = candidate.SkillID
	}
	cases := make([]SkillSelectionCase, 0, len(batch.Cases))
	for _, test := range batch.Cases {
		cases = append(cases, SkillSelectionCase{
			CaseID: test.CaseID,
			Input:  test.Input,
		})
	}
	return SkillSelectionRequest{
		Model:         batch.Model,
		PromptVersion: batch.PromptVersion,
		EvalID:        batch.EvalID + "/run-" + strconv.Itoa(run+1),
		Candidates:    candidates,
		Cases:         cases,
	}, aliasToSkill
}

func validateSkillSelectionOutput(
	raw []byte,
	request SkillSelectionRequest,
) (map[string]struct {
	SelectedCandidate string
	Confidence        float64
}, error) {
	if len(raw) == 0 || len(raw) > maxJudgeInputBytes {
		return nil, errors.New("skill selection output size is invalid")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var output skillSelectionWireOutput
	if err := decoder.Decode(&output); err != nil {
		return nil, fmt.Errorf("decode skill selection output: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("skill selection output contains trailing data")
	}
	if len(output.Results) != len(request.Cases) {
		return nil, errors.New("skill selection output case count does not match request")
	}
	allowedAliases := map[string]struct{}{"none": {}}
	for _, candidate := range request.Candidates {
		allowedAliases[candidate.Alias] = struct{}{}
	}
	expectedCases := make(map[string]struct{}, len(request.Cases))
	for _, test := range request.Cases {
		expectedCases[test.CaseID] = struct{}{}
	}
	results := make(map[string]struct {
		SelectedCandidate string
		Confidence        float64
	}, len(output.Results))
	for _, result := range output.Results {
		if _, found := expectedCases[result.CaseID]; !found {
			return nil, errors.New("skill selection output contains an unknown case")
		}
		if _, duplicate := results[result.CaseID]; duplicate {
			return nil, errors.New("skill selection output duplicates a case")
		}
		if _, found := allowedAliases[result.SelectedCandidate]; !found {
			return nil, errors.New("skill selection output contains an unknown candidate")
		}
		if result.Confidence < 0 || result.Confidence > 1 {
			return nil, errors.New("skill selection confidence is invalid")
		}
		if result.Summary == "" || len(result.Summary) > 500 {
			return nil, errors.New("skill selection summary is invalid")
		}
		results[result.CaseID] = struct {
			SelectedCandidate string
			Confidence        float64
		}{
			SelectedCandidate: result.SelectedCandidate,
			Confidence:        result.Confidence,
		}
	}
	return results, nil
}

func hashSkillEvalCase(test SkillEvalCase) string {
	digest := sha256.Sum256([]byte(test.CaseID + "\x00" + test.Input))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func hasDifferentStrings(values []string) bool {
	if len(values) < 2 {
		return false
	}
	for _, value := range values[1:] {
		if value != values[0] {
			return true
		}
	}
	return false
}
