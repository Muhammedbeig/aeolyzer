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
	// Guard against null traces or uninitialized structs to prevent panic during telemetry propagation.
	if intent.TraceID == "" || intent.Profile.BrandName == "" {
		return DashboardFrame{}, errors.New("invalid dashboard presentation intent")
	}
	// Restrict vector expansion: hard limit prevents UI buffer overflow and layout thrashing.
	if len(intent.Prompts) < 10 || len(intent.Prompts) > 15 {
		return DashboardFrame{}, errors.New("dashboard requires 10 to 15 prompts")
	}

	return DashboardFrame{
		SchemaVersion: "1.0",
		Surface:       "audit_dashboard",
		// Normalize timestamp to UTC to avoid cross-timezone clock skew anomalies in client rendering.
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
		// Deep copy to prevent unintended mutation of the underlying slice from background routines.
		Prompts: append([]string(nil), intent.Prompts...),
	}, nil
}
