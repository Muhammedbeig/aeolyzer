package extensions

import (
	"errors"
	"time"

	"aeolyzer/layer_02_intake"
)

type DashboardIntent struct {
	TraceID     string
	GeneratedAt time.Time
	Profile     intake.ProjectProfile
	Prompts     []string
}

type DashboardFrame struct {
	SchemaVersion string                `json:"schema_version"`
	Surface       string                `json:"surface"`
	GeneratedAt   time.Time             `json:"generated_at"`
	Project       DashboardProject      `json:"project"`
	Tabs          []DashboardNavigation `json:"tabs"`
	Prompts       []string              `json:"prompts"`
}

type DashboardProject struct {
	BrandName string `json:"brand_name"`
	Domain    string `json:"domain"`
	Location  string `json:"location"`
	Language  string `json:"language"`
}

type DashboardNavigation struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
}

func BuildDashboardFrame(intent DashboardIntent) (DashboardFrame, error) {
	if intent.TraceID == "" || intent.Profile.BrandName == "" {
		return DashboardFrame{}, errors.New("invalid dashboard presentation intent")
	}
	if len(intent.Prompts) < 10 || len(intent.Prompts) > 15 {
		return DashboardFrame{}, errors.New("dashboard requires 10 to 15 prompts")
	}

	return DashboardFrame{
		SchemaVersion: "1.0",
		Surface:       "audit_dashboard",
		GeneratedAt:   intent.GeneratedAt.UTC(),
		Project: DashboardProject{
			BrandName: intent.Profile.BrandName,
			Domain:    intent.Profile.Domain,
			Location:  intent.Profile.CountryName,
			Language:  intent.Profile.Language,
		},
		Tabs: []DashboardNavigation{
			{ID: "aeo-insights", Label: "AEO Insights", Enabled: true},
			{ID: "traffic", Label: "Traffic", Enabled: true},
			{ID: "your-prompts", Label: "Your Prompts", Enabled: true},
			{ID: "prompt-research", Label: "Prompt Research", Enabled: true},
			{ID: "site-health", Label: "Site Health", Enabled: true},
		},
		Prompts: append([]string(nil), intent.Prompts...),
	}, nil
}
