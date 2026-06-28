// Package backendconfig loads root policy and routing configuration.
package backendconfig

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

var (
	//go:embed routing-schema.json
	routingJSON []byte
	//go:embed policies.yaml
	policiesYAML []byte
)

// Routing defines the closed intent and mode surface.
type Routing struct {
	Version                     int      `json:"version"`
	PolicyMode                  string   `json:"policy_mode"`
	AllowedIntents              []string `json:"allowed_intents"`
	AllowedModes                []string `json:"allowed_modes"`
	UnknownIntent               string   `json:"unknown_intent"`
	AllowUserSuppliedRoutes     bool     `json:"allow_user_supplied_routes"`
	RequireInternalRouteMapping bool     `json:"require_internal_route_mapping"`
	MinimumConfidence           float64  `json:"minimum_confidence"`
	LowConfidenceBehavior       string   `json:"low_confidence_behavior"`
}

// ActionPolicy is one action-class authorization contract.
type ActionPolicy struct {
	Risk                                       string   `yaml:"risk"`
	AllowedIntents                             []string `yaml:"allowed_intents"`
	RequiresCurrentSourceSafety                bool     `yaml:"requires_current_source_safety,omitempty"`
	RequiresApprovalFor                        string   `yaml:"requires_approval_for,omitempty"`
	RequiresHTTPOrHTTPS                        bool     `yaml:"requires_http_or_https,omitempty"`
	DenyOverwriteWithoutApproval               bool     `yaml:"deny_overwrite_without_approval,omitempty"`
	ReturnSummaryOnly                          bool     `yaml:"return_summary_only,omitempty"`
	RequireAllowedContentType                  bool     `yaml:"require_allowed_content_type,omitempty"`
	RequiresMode                               string   `yaml:"requires_mode,omitempty"`
	DenyDirectOutputPath                       bool     `yaml:"deny_direct_output_path,omitempty"`
	RequiresSelectedText                       bool     `yaml:"requires_selected_text,omitempty"`
	RequireExactMatchPatch                     bool     `yaml:"require_exact_match_patch,omitempty"`
	DenyOverwriteExistingFieldsWithoutApproval bool     `yaml:"deny_overwrite_existing_fields_without_approval,omitempty"`
}

// Policies contains action and semantic denial policy.
type Policies struct {
	Version         int                     `yaml:"version"`
	PolicyMode      string                  `yaml:"policy_mode"`
	ActionClasses   map[string]ActionPolicy `yaml:"action_classes"`
	SemanticDenials map[string]bool         `yaml:"semantic_denials"`
}

// Config contains validated root configuration.
type Config struct {
	Routing  Routing
	Policies Policies
}

// LoadEmbedded strictly parses and validates root configuration.
func LoadEmbedded() (Config, error) {
	var config Config
	decoder := json.NewDecoder(bytes.NewReader(routingJSON))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config.Routing); err != nil {
		return Config{}, fmt.Errorf("decode routing config: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Config{}, errors.New("routing config contains trailing data")
	}
	yamlDecoder := yaml.NewDecoder(bytes.NewReader(policiesYAML))
	yamlDecoder.KnownFields(true)
	if err := yamlDecoder.Decode(&config.Policies); err != nil {
		return Config{}, fmt.Errorf("decode action policies: %w", err)
	}
	if err := yamlDecoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Config{}, errors.New("action policies contain multiple documents")
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

// Validate checks fail-closed cross-references and action invariants.
func (c Config) Validate() error {
	if c.Routing.Version != 2 ||
		c.Routing.PolicyMode != "fail_closed" ||
		len(c.Routing.AllowedIntents) == 0 ||
		len(c.Routing.AllowedModes) == 0 ||
		c.Routing.UnknownIntent != "reject" ||
		c.Routing.AllowUserSuppliedRoutes ||
		!c.Routing.RequireInternalRouteMapping ||
		c.Routing.MinimumConfidence <= 0 ||
		c.Routing.MinimumConfidence > 1 ||
		c.Routing.LowConfidenceBehavior != "clarify" {
		return errors.New("routing configuration is unsafe")
	}
	if c.Policies.Version != 2 ||
		c.Policies.PolicyMode != "fail_closed" ||
		len(c.Policies.ActionClasses) == 0 ||
		len(c.Policies.SemanticDenials) == 0 {
		return errors.New("action policy configuration is incomplete")
	}
	intents := make(map[string]struct{}, len(c.Routing.AllowedIntents))
	for _, intent := range c.Routing.AllowedIntents {
		if intent == "" {
			return errors.New("routing configuration contains empty intent")
		}
		intents[intent] = struct{}{}
	}
	for class, policy := range c.Policies.ActionClasses {
		if class == "" || len(policy.AllowedIntents) == 0 {
			return errors.New("action policy contains empty class or intents")
		}
		switch policy.Risk {
		case "low", "medium", "high", "critical":
		default:
			return errors.New("action policy contains invalid risk")
		}
		for _, intent := range policy.AllowedIntents {
			if _, found := intents[intent]; !found {
				return fmt.Errorf("action class %s references unknown intent %s", class, intent)
			}
		}
	}
	for denial, enabled := range c.Policies.SemanticDenials {
		if denial == "" || !enabled {
			return errors.New("semantic denial policy must fail closed")
		}
	}
	return nil
}
