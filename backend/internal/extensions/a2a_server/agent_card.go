// Package a2aserver validates public A2A application contracts.
package a2aserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"aeolyzer/internal/extensions"

	"github.com/a2aproject/a2a-go/v2/a2a"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/server/adka2a/v2"
)

const (
	defaultA2APath          = "/a2a"
	defaultAgentVersion     = "1.0.0"
	defaultSecurityScheme   = "googleOidc"
	defaultOpenIDConnectURL = "https://accounts.google.com/.well-known/openid-configuration"
)

// AgentCard is the canonical A2A Agent Card type.
type AgentCard = a2a.AgentCard

// AgentCardConfig contains public, non-secret Agent Card inputs.
type AgentCardConfig struct {
	Name               string
	Description        string
	PublicBaseURL      string
	A2APath            string
	Version            string
	ProviderName       string
	ProviderURL        string
	DocumentationURL   string
	IconURL            string
	SecuritySchemeName string
	OpenIDConnectURL   string
	Skills             []a2a.AgentSkill
	InputModes         []string
	OutputModes        []string
}

// DefaultPublicSkills returns the product-level skills safe for public A2A discovery.
func DefaultPublicSkills() []a2a.AgentSkill {
	return []a2a.AgentSkill{{
		ID:          "site_visibility_guidance",
		Name:        "Site visibility guidance",
		Description: "Explains safe, public website visibility and content improvement options without exposing internal topology.",
		Tags:        []string{"aeo", "seo", "content"},
	}}
}

// BuildAgentSkills converts ADK agent metadata into canonical A2A skills.
func BuildAgentSkills(adkAgent agent.Agent) []a2a.AgentSkill {
	if adkAgent == nil {
		return nil
	}
	return adka2a.BuildAgentSkills(adkAgent)
}

// NewAgentCard builds a canonical, schema-valid A2A Agent Card.
func NewAgentCard(config AgentCardConfig) (*a2a.AgentCard, error) {
	cardURL, err := agentEndpoint(config.PublicBaseURL, config.A2APath)
	if err != nil {
		return nil, err
	}
	skills := normalizeSkills(config.Skills)
	if len(skills) == 0 {
		return nil, errors.New("agent card requires at least one public skill")
	}
	version := strings.TrimSpace(config.Version)
	if version == "" {
		version = defaultAgentVersion
	}
	name := strings.TrimSpace(config.Name)
	description := strings.TrimSpace(config.Description)
	if name == "" || description == "" {
		return nil, errors.New("agent card name and description are required")
	}
	inputModes := defaultModes(config.InputModes)
	outputModes := defaultModes(config.OutputModes)

	card := &a2a.AgentCard{
		Name:        name,
		Description: description,
		SupportedInterfaces: []*a2a.AgentInterface{
			a2a.NewAgentInterface(cardURL, a2a.TransportProtocolJSONRPC),
		},
		Capabilities:       a2a.AgentCapabilities{},
		DefaultInputModes:  inputModes,
		DefaultOutputModes: outputModes,
		DocumentationURL:   strings.TrimSpace(config.DocumentationURL),
		IconURL:            strings.TrimSpace(config.IconURL),
		Skills:             skills,
		Version:            version,
	}
	if config.ProviderName != "" || config.ProviderURL != "" {
		card.Provider = &a2a.AgentProvider{
			Org: strings.TrimSpace(config.ProviderName),
			URL: strings.TrimSpace(config.ProviderURL),
		}
	}
	addSecurity(card, config.SecuritySchemeName, config.OpenIDConnectURL)
	if err := ValidateAgentCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

// ValidateAgentCard validates canonical A2A schema and blocks protected topology disclosure.
func ValidateAgentCard(card *a2a.AgentCard) error {
	if card == nil {
		return errors.New("agent card is required")
	}
	if err := validateAgentCardSemantics(card); err != nil {
		return err
	}
	schemas, err := extensions.NewSchemas()
	if err != nil {
		return err
	}
	data, err := json.Marshal(card)
	if err != nil {
		return errors.New("agent card cannot be encoded")
	}
	if err := schemas.ValidateJSON(extensions.ContractA2AAgentCard, data); err != nil {
		return err
	}
	if containsProtectedMetadata(card) {
		return errors.New("agent card discloses protected metadata")
	}
	return nil
}

func agentEndpoint(baseURL, a2aPath string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil ||
		parsed.Scheme != "https" ||
		parsed.Hostname() == "" ||
		parsed.User != nil ||
		parsed.Fragment != "" {
		return "", errors.New("agent card public base url is invalid")
	}
	cleanPath := strings.TrimSpace(a2aPath)
	if cleanPath == "" {
		cleanPath = defaultA2APath
	}
	if !strings.HasPrefix(cleanPath, "/") || strings.Contains(cleanPath, "..") {
		return "", errors.New("agent card a2a path is invalid")
	}
	parsed.Path = path.Join(parsed.Path, cleanPath)
	parsed.RawQuery = ""
	return parsed.String(), nil
}

func normalizeSkills(skills []a2a.AgentSkill) []a2a.AgentSkill {
	normalized := make([]a2a.AgentSkill, 0, len(skills))
	for _, skill := range skills {
		skill.ID = strings.TrimSpace(skill.ID)
		skill.Name = strings.TrimSpace(skill.Name)
		skill.Description = strings.TrimSpace(skill.Description)
		if len(skill.Tags) == 0 {
			skill.Tags = []string{"aeo"}
		}
		normalized = append(normalized, skill)
	}
	return normalized
}

func defaultModes(modes []string) []string {
	if len(modes) == 0 {
		return []string{"text/plain"}
	}
	normalized := make([]string, 0, len(modes))
	for _, mode := range modes {
		if trimmed := strings.TrimSpace(mode); trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	if len(normalized) == 0 {
		return []string{"text/plain"}
	}
	return normalized
}

func addSecurity(card *a2a.AgentCard, schemeName, openIDConnectURL string) {
	name := strings.TrimSpace(schemeName)
	if name == "" {
		name = defaultSecurityScheme
	}
	discoveryURL := strings.TrimSpace(openIDConnectURL)
	if discoveryURL == "" {
		discoveryURL = defaultOpenIDConnectURL
	}
	card.SecuritySchemes = a2a.NamedSecuritySchemes{
		a2a.SecuritySchemeName(name): a2a.OpenIDConnectSecurityScheme{
			OpenIDConnectURL: discoveryURL,
		},
	}
	card.SecurityRequirements = a2a.SecurityRequirementsOptions{
		a2a.SecurityRequirements{
			a2a.SecuritySchemeName(name): a2a.SecuritySchemeScopes{},
		},
	}
}

func validateAgentCardSemantics(card *a2a.AgentCard) error {
	if strings.TrimSpace(card.Name) == "" || strings.TrimSpace(card.Description) == "" {
		return errors.New("agent card name and description are required")
	}
	if card.Version == "" || len(card.Skills) == 0 || len(card.SupportedInterfaces) == 0 {
		return errors.New("agent card required fields are missing")
	}
	for _, iface := range card.SupportedInterfaces {
		if iface == nil ||
			iface.ProtocolVersion != a2a.Version ||
			iface.ProtocolBinding != a2a.TransportProtocolJSONRPC ||
			!isPublicHTTPS(iface.URL) {
			return errors.New("agent card supported interface is invalid")
		}
	}
	for _, skill := range card.Skills {
		if skill.ID == "" || skill.Name == "" || skill.Description == "" || len(skill.Tags) == 0 {
			return errors.New("agent card public skill is invalid")
		}
	}
	return nil
}

func isPublicHTTPS(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil &&
		parsed.Scheme == "https" &&
		parsed.Hostname() != "" &&
		parsed.User == nil &&
		parsed.Fragment == ""
}

func containsProtectedMetadata(card *a2a.AgentCard) bool {
	data, err := json.Marshal(card)
	if err != nil {
		return true
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return true
	}
	return scanStrings(decoded)
}

func scanStrings(value any) bool {
	switch v := value.(type) {
	case string:
		return containsForbiddenFragment(v)
	case []any:
		for _, item := range v {
			if scanStrings(item) {
				return true
			}
		}
	case map[string]any:
		for key, item := range v {
			if containsForbiddenKey(key) || scanStrings(item) {
				return true
			}
		}
	}
	return false
}

func containsForbiddenKey(key string) bool {
	switch strings.ToLower(key) {
	case "workflow_id", "profile_id", "trace_id", "mcp_endpoint", "tool_id", "skill_path":
		return true
	default:
		return false
	}
}

func containsForbiddenFragment(value string) bool {
	lower := strings.ToLower(value)
	for _, forbidden := range []string{
		"internal_",
		"/internal/",
		"workflow id",
		"workflow_id",
		"profile id",
		"profile_id",
		"trace id",
		"trace_id",
		"mcp://",
		"mcp endpoint",
		"mcp server",
		"sandbox",
		"jit token",
		"credential broker",
		"skill path",
		"skill file",
		"skill.md",
		"policy file",
		"sql executor",
	} {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	return false
}

func agentCardJSON(card *a2a.AgentCard) ([]byte, error) {
	if err := ValidateAgentCard(card); err != nil {
		return nil, err
	}
	data, err := json.Marshal(card)
	if err != nil {
		return nil, fmt.Errorf("marshal agent card: %w", err)
	}
	return data, nil
}
