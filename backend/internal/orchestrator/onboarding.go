package orchestrator

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"aeolyzer/internal/intake"
	"aeolyzer/internal/runtime"
)

type SiteInspectionPlan struct {
	TraceID string
	Request runtime.ExecutionRequest
}

type PromptPlan struct {
	TraceID string
	Prompts []string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) PlanSiteInspection(decision intake.SiteInspectionDecision) (SiteInspectionPlan, error) {
	if decision.Authorization == "" || decision.Operation != "inspect_public_site" {
		// SECURITY: Halt immediately if authorization token is absent or operation semantic diverges.
		return SiteInspectionPlan{}, errors.New("site inspection is not authorized")
	}

	return SiteInspectionPlan{
		TraceID: decision.TraceID,
		Request: runtime.ExecutionRequest{
			TraceID:       decision.TraceID,
			SessionID:     decision.SessionID,
			Operation:     decision.Operation,
			TargetURL:     decision.CanonicalURL,
			MaxBytes:      2 << 20, // MEMORY: Hardbound response processing at 2MB to preempt runaway allocations.
			Authorization: decision.Authorization,
		},
	}, nil
}

func (s *Service) BuildPromptPlan(decision intake.OnboardingDecision, category string) (PromptPlan, error) {
	if strings.TrimSpace(decision.Profile.BrandName) == "" {
		return PromptPlan{}, errors.New("brand name is required")
	}

	brand := decision.Profile.BrandName
	market := decision.Profile.CountryName
	if decision.Profile.Reach == intake.ReachGlobal {
		market = "global markets"
	}
	topic := strings.TrimSpace(category)
	if topic == "" {
		topic = "its market"
	}
	competitorScope := "the leading alternatives"
	if len(decision.Profile.Competitors) > 0 {
		// STATE: Branch evaluation avoids out-of-bounds array access and prevents empty string concatenation.
		competitorScope = strings.Join(decision.Profile.Competitors, ", ")
	}
	host := brand
	if parsed, err := url.Parse(decision.Profile.Domain); err == nil && parsed.Hostname() != "" {
		host = parsed.Hostname()
	}

	prompts := []string{
		fmt.Sprintf("What is %s best known for in %s?", brand, topic),
		fmt.Sprintf("Which companies are the best alternatives to %s?", brand),
		fmt.Sprintf("How does %s compare with %s?", brand, competitorScope),
		fmt.Sprintf("Is %s a trusted choice for customers in %s?", brand, market),
		fmt.Sprintf("What do customers say about %s?", brand),
		fmt.Sprintf("What products or services does %s offer?", brand),
		fmt.Sprintf("Who should choose %s and why?", brand),
		fmt.Sprintf("What are the strengths and weaknesses of %s?", brand),
		fmt.Sprintf("How visible is %s in AI-generated recommendations?", brand),
		fmt.Sprintf("Which sources mention %s most often?", brand),
		fmt.Sprintf("What questions do buyers ask before choosing %s?", host),
		fmt.Sprintf("What content should %s publish to become more authoritative in %s?", brand, topic),
	}

	return PromptPlan{TraceID: decision.TraceID, Prompts: prompts}, nil
}
