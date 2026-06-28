// Command skilleval runs the embedded Layer 4 trigger corpus against a real,
// environment-authenticated model and emits sanitized Layer 8 evidence.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"

	llmprovider "aeolyzer/internal/interop/llm_provider"
	observabilityconfig "aeolyzer/internal/observability/config"
	evaluationengine "aeolyzer/internal/observability/evaluation_engine"
	"aeolyzer/internal/skills"
)

const (
	defaultGeminiModel        = "gemini-3.5-flash"
	defaultSkillPromptVersion = "skill-router-eval-v1"
	defaultCandidateCount     = 10
)

type options struct {
	model         string
	promptVersion string
	skillID       string
	output        string
}

type suiteReport struct {
	CorpusChecksum  string                             `json:"corpus_checksum"`
	Model           string                             `json:"model"`
	PromptVersion   string                             `json:"prompt_version"`
	SkillsEvaluated int                                `json:"skills_evaluated"`
	SkillsPassed    int                                `json:"skills_passed"`
	SkillsFailed    int                                `json:"skills_failed"`
	Passed          bool                               `json:"passed"`
	Reports         []evaluationengine.SkillEvalReport `json:"reports"`
}

func main() {
	var config options
	flag.StringVar(&config.model, "model", defaultGeminiModel, "pinned Gemini model")
	flag.StringVar(
		&config.promptVersion,
		"prompt-version",
		defaultSkillPromptVersion,
		"versioned skill-routing evaluation prompt",
	)
	flag.StringVar(&config.skillID, "skill", "", "evaluate only one skill id")
	flag.StringVar(&config.output, "output", "", "write sanitized JSON evidence to this file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	if err := run(ctx, config); err != nil {
		fmt.Fprintln(os.Stderr, "skill evaluation failed:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, config options) error {
	if config.model == "" || config.promptVersion == "" {
		return errors.New("model and prompt version are required")
	}
	corpus, err := skills.LoadEmbeddedTriggerEvalCorpus()
	if err != nil {
		return fmt.Errorf("load embedded trigger corpus: %w", err)
	}
	policies, err := observabilityconfig.LoadEmbeddedPolicies()
	if err != nil {
		return fmt.Errorf("load eval policy: %w", err)
	}
	batches, err := buildSkillEvalBatches(
		corpus,
		config.model,
		config.promptVersion,
		config.skillID,
	)
	if err != nil {
		return err
	}
	client, err := llmprovider.NewGeminiJudgeClientFromEnvironment(http.DefaultTransport)
	if err != nil {
		return err
	}

	report := suiteReport{
		CorpusChecksum: corpus.Checksum,
		Model:          config.model,
		PromptVersion:  config.promptVersion,
		Passed:         true,
		Reports:        make([]evaluationengine.SkillEvalReport, 0, len(batches)),
	}
	for _, batch := range batches {
		passK, err := passKForTier(batch.tier, policies.Eval)
		if err != nil {
			return err
		}
		result, err := evaluationengine.EvaluateSkillTriggers(
			ctx,
			client,
			batch.batch,
			passK,
			policies.Eval.Judge.MinimumConfidence,
		)
		if err != nil {
			return fmt.Errorf("evaluate skill %q: %w", batch.batch.TargetSkillID, err)
		}
		report.Reports = append(report.Reports, result)
		report.SkillsEvaluated++
		if result.Passed {
			report.SkillsPassed++
		} else {
			report.SkillsFailed++
			report.Passed = false
		}
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode skill evaluation report: %w", err)
	}
	data = append(data, '\n')
	if config.output != "" {
		if err := os.WriteFile(config.output, data, 0o600); err != nil {
			return fmt.Errorf("write skill evaluation report: %w", err)
		}
	} else {
		if _, err := os.Stdout.Write(data); err != nil {
			return fmt.Errorf("write skill evaluation report: %w", err)
		}
	}
	if !report.Passed {
		return errors.New("one or more skill evaluation gates failed")
	}
	return nil
}

type tieredBatch struct {
	tier  string
	batch evaluationengine.SkillEvalBatch
}

func buildSkillEvalBatches(
	corpus skills.TriggerEvalCorpus,
	model string,
	promptVersion string,
	filterSkillID string,
) ([]tieredBatch, error) {
	skillByID := make(map[string]skills.EvalSkill, len(corpus.Skills))
	for _, skill := range corpus.Skills {
		skillByID[skill.SkillID] = skill
	}
	if filterSkillID != "" {
		if _, found := skillByID[filterSkillID]; !found {
			return nil, fmt.Errorf("skill %q is not registered", filterSkillID)
		}
	}

	var batches []tieredBatch
	for _, target := range corpus.Skills {
		if filterSkillID != "" && target.SkillID != filterSkillID {
			continue
		}
		var targetCases []skills.TriggerEvalCase
		for _, test := range corpus.Cases {
			if test.TargetSkillID == target.SkillID {
				targetCases = append(targetCases, test)
			}
		}
		candidateIDs := selectCandidateSkills(
			target,
			targetCases,
			corpus.Skills,
			defaultCandidateCount,
		)
		candidates := make([]evaluationengine.SkillEvalCandidate, 0, len(candidateIDs))
		for _, skillID := range candidateIDs {
			candidate := skillByID[skillID]
			candidates = append(candidates, evaluationengine.SkillEvalCandidate{
				SkillID:      candidate.SkillID,
				Description:  candidate.Description,
				AntiTriggers: append([]string(nil), candidate.AntiTriggers...),
			})
		}
		cases := make([]evaluationengine.SkillEvalCase, 0, len(targetCases))
		for _, test := range targetCases {
			cases = append(cases, evaluationengine.SkillEvalCase{
				CaseID:         test.CaseID,
				Group:          test.Group,
				Input:          test.Input,
				ExpectedSkill:  test.ExpectedSkill,
				ForbiddenSkill: test.ForbiddenSkill,
				SafetyCritical: test.SafetyCritical,
			})
		}
		batches = append(batches, tieredBatch{
			tier: target.Tier,
			batch: evaluationengine.SkillEvalBatch{
				EvalID:         "skill-trigger/" + target.SkillID,
				CorpusChecksum: corpus.Checksum,
				TargetSkillID:  target.SkillID,
				Model:          model,
				PromptVersion:  promptVersion,
				Candidates:     candidates,
				Cases:          cases,
			},
		})
	}
	return batches, nil
}

func selectCandidateSkills(
	target skills.EvalSkill,
	cases []skills.TriggerEvalCase,
	catalog []skills.EvalSkill,
	limit int,
) []string {
	selected := make(map[string]struct{}, limit)
	result := make([]string, 0, limit)
	add := func(skillID string) {
		if skillID == "" || len(result) >= limit {
			return
		}
		if _, found := selected[skillID]; found {
			return
		}
		selected[skillID] = struct{}{}
		result = append(result, skillID)
	}
	add(target.SkillID)
	for _, test := range cases {
		add(test.ExpectedSkill)
		add(test.ForbiddenSkill)
	}

	type scoredSkill struct {
		skillID string
		score   int
	}
	var scored []scoredSkill
	for _, candidate := range catalog {
		if _, found := selected[candidate.SkillID]; found {
			continue
		}
		scored = append(scored, scoredSkill{
			skillID: candidate.SkillID,
			score: intersectionCount(
				target.CompatibleIntents,
				candidate.CompatibleIntents,
			)*2 + intersectionCount(target.CapabilityTags, candidate.CapabilityTags),
		})
	}
	sort.SliceStable(scored, func(left, right int) bool {
		if scored[left].score == scored[right].score {
			return scored[left].skillID < scored[right].skillID
		}
		return scored[left].score > scored[right].score
	})
	for _, candidate := range scored {
		add(candidate.skillID)
	}
	return result
}

func intersectionCount(left, right []string) int {
	values := make(map[string]struct{}, len(left))
	for _, value := range left {
		values[value] = struct{}{}
	}
	count := 0
	for _, value := range right {
		if _, found := values[value]; found {
			count++
		}
	}
	return count
}

func passKForTier(
	tier string,
	policy observabilityconfig.EvalPolicy,
) (int, error) {
	switch tier {
	case "read":
		return policy.PassK.ReadWorkflow, nil
	case "draft":
		return policy.PassK.DraftWorkflow, nil
	case "act":
		return policy.PassK.GuardedWriteWorkflow, nil
	default:
		return 0, fmt.Errorf("skill tier %q has no pass k policy", tier)
	}
}
