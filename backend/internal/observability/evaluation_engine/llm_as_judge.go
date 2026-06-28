package evaluationengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
)

const (
	maxJudgeInputBytes = 256 << 10
	maxJudgeRetries    = 2
)

// PairwiseWinner is the judge's pairwise preference.
type PairwiseWinner string

const (
	// WinnerA means candidate A better satisfies the rubric.
	WinnerA PairwiseWinner = "A"
	// WinnerB means candidate B better satisfies the rubric.
	WinnerB PairwiseWinner = "B"
	// WinnerTie means neither candidate is materially better.
	WinnerTie PairwiseWinner = "TIE"
)

// RubricDimension defines one scored dimension.
type RubricDimension struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Minimum     int     `json:"minimum"`
	Maximum     int     `json:"maximum"`
}

// Rubric is the complete, versioned scoring contract.
type Rubric struct {
	ID         string            `json:"id"`
	Version    string            `json:"version"`
	Dimensions []RubricDimension `json:"dimensions"`
	PassScore  float64           `json:"pass_score"`
}

// JudgeInput contains redacted evaluation material.
type JudgeInput struct {
	EvalID      string `json:"eval_id"`
	Instruction string `json:"instruction"`
	CandidateA  string `json:"candidate_a"`
	CandidateB  string `json:"candidate_b,omitempty"`
	Pairwise    bool   `json:"pairwise"`
	Redacted    bool   `json:"redacted"`
}

// JudgeConfig pins the non-deterministic judge boundary.
type JudgeConfig struct {
	Model                 string  `json:"model"`
	PromptVersion         string  `json:"prompt_version"`
	Temperature           float64 `json:"temperature"`
	MinimumConfidence     float64 `json:"minimum_confidence"`
	HumanReviewConfidence float64 `json:"human_review_confidence"`
	MaxRetries            int     `json:"max_retries"`
}

// JudgeRequest is the provider-neutral request sent to a judge client.
type JudgeRequest struct {
	Model         string     `json:"model"`
	PromptVersion string     `json:"prompt_version"`
	Temperature   float64    `json:"temperature"`
	Input         JudgeInput `json:"input"`
	Rubric        Rubric     `json:"rubric"`
}

// JudgeClient calls one configured judge provider. Implementations must not log
// request payloads or credentials.
type JudgeClient interface {
	Judge(context.Context, JudgeRequest) ([]byte, error)
}

// ValidateJudgeRequest validates the provider-neutral request contract before
// a Layer 7 provider adapter transmits it. Providers must still enforce their
// own transport, credential, response-size, and redirect controls.
func ValidateJudgeRequest(request JudgeRequest) error {
	if request.Model == "" || request.PromptVersion == "" {
		return errors.New("judge model and prompt version are required")
	}
	if request.Temperature != 0 {
		return errors.New("judge temperature must be zero")
	}
	if err := validateJudgeInput(request.Input); err != nil {
		return err
	}
	return validateRubric(request.Rubric)
}

// JudgeScore is the validated and deterministically recomputed judge result.
type JudgeScore struct {
	RubricID      string         `json:"rubric_id"`
	Dimensions    map[string]int `json:"dimensions"`
	WeightedScore float64        `json:"weighted_score"`
	Pass          bool           `json:"pass"`
	Winner        PairwiseWinner `json:"winner,omitempty"`
	Confidence    float64        `json:"confidence"`
	Summary       string         `json:"summary"`
	RequiresHuman bool           `json:"requires_human_review"`
}

// BiasControlledScore contains both position-swapped runs.
type BiasControlledScore struct {
	Forward       JudgeScore `json:"forward"`
	Swapped       JudgeScore `json:"swapped"`
	Stable        bool       `json:"stable"`
	Canonical     JudgeScore `json:"canonical"`
	RequiresHuman bool       `json:"requires_human_review"`
}

type judgeWireOutput struct {
	RubricID   string         `json:"rubric_id"`
	Dimensions map[string]int `json:"dimensions"`
	Winner     PairwiseWinner `json:"winner,omitempty"`
	Confidence float64        `json:"confidence"`
	Summary    string         `json:"summary"`
}

// ScoreWithJudge calls a pinned judge and validates its result against the
// deterministic rubric contract.
func ScoreWithJudge(
	ctx context.Context,
	client JudgeClient,
	input JudgeInput,
	rubric Rubric,
	config JudgeConfig,
) (JudgeScore, error) {
	if client == nil {
		return JudgeScore{}, errors.New("judge client is required")
	}
	if err := validateJudgeInput(input); err != nil {
		return JudgeScore{}, err
	}
	if err := validateRubric(rubric); err != nil {
		return JudgeScore{}, err
	}
	if err := validateJudgeConfig(config); err != nil {
		return JudgeScore{}, err
	}

	request := JudgeRequest{
		Model:         config.Model,
		PromptVersion: config.PromptVersion,
		Temperature:   config.Temperature,
		Input:         input,
		Rubric:        rubric,
	}

	var lastErr error
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return JudgeScore{}, fmt.Errorf("score with judge: %w", err)
		}
		raw, err := client.Judge(ctx, request)
		if err != nil {
			lastErr = fmt.Errorf("call judge: %w", err)
			continue
		}
		score, err := validateJudgeOutput(raw, input, rubric, config)
		if err == nil {
			return score, nil
		}
		lastErr = err
	}
	return JudgeScore{}, fmt.Errorf("judge output rejected after %d attempts: %w", config.MaxRetries+1, lastErr)
}

// RunPositionSwap evaluates pairwise candidates in both positions and detects
// ordering bias.
func RunPositionSwap(
	ctx context.Context,
	client JudgeClient,
	input JudgeInput,
	rubric Rubric,
	config JudgeConfig,
) (BiasControlledScore, error) {
	if !input.Pairwise || input.CandidateB == "" {
		return BiasControlledScore{}, errors.New("position swap requires two pairwise candidates")
	}

	forward, err := ScoreWithJudge(ctx, client, input, rubric, config)
	if err != nil {
		return BiasControlledScore{}, err
	}
	swappedInput := input
	swappedInput.CandidateA = input.CandidateB
	swappedInput.CandidateB = input.CandidateA
	swapped, err := ScoreWithJudge(ctx, client, swappedInput, rubric, config)
	if err != nil {
		return BiasControlledScore{}, err
	}

	canonicalSwappedWinner := invertWinner(swapped.Winner)
	stable := forward.Winner == canonicalSwappedWinner &&
		math.Abs(forward.WeightedScore-swapped.WeightedScore) <= 0.5
	requiresHuman := !stable || forward.RequiresHuman || swapped.RequiresHuman
	canonical := forward
	canonical.RequiresHuman = requiresHuman
	if !stable {
		canonical.Pass = false
	}
	return BiasControlledScore{
		Forward:       forward,
		Swapped:       swapped,
		Stable:        stable,
		Canonical:     canonical,
		RequiresHuman: requiresHuman,
	}, nil
}

// ValidateJudgeOutput validates strict provider output independent of a rubric.
// Rubric-specific dimension and weighting checks occur in ScoreWithJudge.
func ValidateJudgeOutput(raw []byte) (JudgeScore, error) {
	output, err := decodeJudgeOutput(raw)
	if err != nil {
		return JudgeScore{}, err
	}
	if err := validateWireOutput(output); err != nil {
		return JudgeScore{}, err
	}
	return JudgeScore{
		RubricID:   output.RubricID,
		Dimensions: output.Dimensions,
		Winner:     output.Winner,
		Confidence: output.Confidence,
		Summary:    output.Summary,
	}, nil
}

// DetectJudgeInstability reports material score, winner, or pass disagreement.
func DetectJudgeInstability(scores []JudgeScore) bool {
	if len(scores) < 2 {
		return false
	}
	baseline := scores[0]
	for _, score := range scores[1:] {
		if score.Winner != baseline.Winner ||
			score.Pass != baseline.Pass ||
			math.Abs(score.WeightedScore-baseline.WeightedScore) > 0.5 {
			return true
		}
	}
	return false
}

func validateJudgeOutput(
	raw []byte,
	input JudgeInput,
	rubric Rubric,
	config JudgeConfig,
) (JudgeScore, error) {
	output, err := decodeJudgeOutput(raw)
	if err != nil {
		return JudgeScore{}, err
	}
	if err := validateWireOutput(output); err != nil {
		return JudgeScore{}, err
	}
	if output.RubricID != rubric.ID {
		return JudgeScore{}, errors.New("judge returned the wrong rubric id")
	}
	if input.Pairwise && output.Winner == "" {
		return JudgeScore{}, errors.New("pairwise judge result requires a winner")
	}
	if !input.Pairwise && output.Winner != "" {
		return JudgeScore{}, errors.New("single-output judge result must not include a winner")
	}

	weightedScore, err := calculateWeightedScore(output.Dimensions, rubric)
	if err != nil {
		return JudgeScore{}, err
	}
	return JudgeScore{
		RubricID:      output.RubricID,
		Dimensions:    output.Dimensions,
		WeightedScore: weightedScore,
		Pass:          weightedScore >= rubric.PassScore && output.Confidence >= config.MinimumConfidence,
		Winner:        output.Winner,
		Confidence:    output.Confidence,
		Summary:       output.Summary,
		RequiresHuman: output.Confidence < config.HumanReviewConfidence,
	}, nil
}

func decodeJudgeOutput(raw []byte) (judgeWireOutput, error) {
	if len(raw) == 0 {
		return judgeWireOutput{}, errors.New("judge output is empty")
	}
	if len(raw) > maxJudgeInputBytes {
		return judgeWireOutput{}, errors.New("judge output exceeds size limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var output judgeWireOutput
	if err := decoder.Decode(&output); err != nil {
		return judgeWireOutput{}, fmt.Errorf("decode judge json: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return judgeWireOutput{}, errors.New("judge output contains trailing data")
	}
	return output, nil
}

func validateWireOutput(output judgeWireOutput) error {
	if output.RubricID == "" {
		return errors.New("judge rubric id is required")
	}
	if len(output.Dimensions) == 0 {
		return errors.New("judge dimensions are required")
	}
	if output.Confidence < 0 || output.Confidence > 1 {
		return errors.New("judge confidence must be between zero and one")
	}
	if output.Summary == "" || len(output.Summary) > 1000 {
		return errors.New("judge summary length is invalid")
	}
	switch output.Winner {
	case "", WinnerA, WinnerB, WinnerTie:
	default:
		return errors.New("judge winner is invalid")
	}
	for id, score := range output.Dimensions {
		if id == "" || score < 1 || score > 5 {
			return errors.New("judge dimension score is invalid")
		}
	}
	return nil
}

func validateJudgeInput(input JudgeInput) error {
	if input.EvalID == "" {
		return errors.New("judge eval id is required")
	}
	if !input.Redacted {
		return errors.New("judge input must be redacted")
	}
	if input.Instruction == "" || input.CandidateA == "" {
		return errors.New("judge instruction and candidate a are required")
	}
	if input.Pairwise != (input.CandidateB != "") {
		return errors.New("pairwise judge input requires exactly two candidates")
	}
	size := len(input.Instruction) + len(input.CandidateA) + len(input.CandidateB)
	if size > maxJudgeInputBytes {
		return errors.New("judge input exceeds size limit")
	}
	return nil
}

func validateRubric(rubric Rubric) error {
	if rubric.ID == "" || rubric.Version == "" {
		return errors.New("rubric id and version are required")
	}
	if len(rubric.Dimensions) == 0 || len(rubric.Dimensions) > 32 {
		return errors.New("rubric must contain between one and 32 dimensions")
	}
	if rubric.PassScore < 1 || rubric.PassScore > 5 {
		return errors.New("rubric pass score must be between one and five")
	}
	seen := make(map[string]struct{}, len(rubric.Dimensions))
	totalWeight := 0.0
	for _, dimension := range rubric.Dimensions {
		if dimension.ID == "" || dimension.Description == "" {
			return errors.New("rubric dimension id and description are required")
		}
		if _, exists := seen[dimension.ID]; exists {
			return fmt.Errorf("duplicate rubric dimension %q", dimension.ID)
		}
		seen[dimension.ID] = struct{}{}
		if dimension.Weight <= 0 || dimension.Weight > 1 {
			return fmt.Errorf("rubric dimension %q has invalid weight", dimension.ID)
		}
		if dimension.Minimum < 1 || dimension.Maximum > 5 || dimension.Minimum > dimension.Maximum {
			return fmt.Errorf("rubric dimension %q has invalid bounds", dimension.ID)
		}
		totalWeight += dimension.Weight
	}
	if math.Abs(totalWeight-1) > 0.000001 {
		return fmt.Errorf("rubric weights total %f, want 1", totalWeight)
	}
	return nil
}

func validateJudgeConfig(config JudgeConfig) error {
	if config.Model == "" || config.PromptVersion == "" {
		return errors.New("judge model and prompt version are required")
	}
	if config.Temperature != 0 {
		return errors.New("judge temperature must be zero")
	}
	if config.MinimumConfidence < 0 || config.MinimumConfidence > 1 {
		return errors.New("minimum confidence must be between zero and one")
	}
	if config.HumanReviewConfidence < config.MinimumConfidence ||
		config.HumanReviewConfidence > 1 {
		return errors.New("human review confidence threshold is invalid")
	}
	if config.MaxRetries < 0 || config.MaxRetries > maxJudgeRetries {
		return fmt.Errorf("judge max retries must be between zero and %d", maxJudgeRetries)
	}
	return nil
}

func calculateWeightedScore(scores map[string]int, rubric Rubric) (float64, error) {
	if len(scores) != len(rubric.Dimensions) {
		return 0, errors.New("judge returned an incomplete dimension set")
	}
	dimensions := append([]RubricDimension(nil), rubric.Dimensions...)
	sort.Slice(dimensions, func(i, j int) bool {
		return dimensions[i].ID < dimensions[j].ID
	})

	total := 0.0
	for _, dimension := range dimensions {
		score, found := scores[dimension.ID]
		if !found {
			return 0, fmt.Errorf("judge omitted dimension %q", dimension.ID)
		}
		if score < dimension.Minimum || score > dimension.Maximum {
			return 0, fmt.Errorf("judge score for %q is outside rubric bounds", dimension.ID)
		}
		total += float64(score) * dimension.Weight
	}
	return total, nil
}

func invertWinner(winner PairwiseWinner) PairwiseWinner {
	switch winner {
	case WinnerA:
		return WinnerB
	case WinnerB:
		return WinnerA
	default:
		return winner
	}
}
