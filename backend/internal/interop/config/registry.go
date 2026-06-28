// Package interopconfig loads and validates Layer 7 connector configuration.
package interopconfig

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"

	"aeolyzer/internal/schemavalidation"
	"go.yaml.in/yaml/v3"
)

var (
	//go:embed connector-registry.yaml
	connectorRegistryYAML []byte
	//go:embed source-contracts.yaml
	sourceContractsYAML []byte
	//go:embed mcp-server-manifest.schema.json
	mcpServerManifestSchemaJSON []byte
)

// Registry is the fail-closed connector allowlist.
type Registry struct {
	Version       int         `yaml:"version"`
	PolicyMode    string      `yaml:"policy_mode"`
	DefaultStatus string      `yaml:"default_status"`
	Connectors    []Connector `yaml:"connectors"`
}

// Connector is one bounded connector registration.
type Connector struct {
	ConnectorID           string   `yaml:"connector_id"`
	Status                string   `yaml:"status"`
	OwnerTeam             string   `yaml:"owner_team"`
	SourceType            string   `yaml:"source_type"`
	Transport             string   `yaml:"transport"`
	Authentication        string   `yaml:"authentication"`
	TenantIsolation       string   `yaml:"tenant_isolation"`
	ReadOnly              bool     `yaml:"read_only"`
	AllowedActionClasses  []string `yaml:"allowed_action_classes"`
	RequiredRuntimeClass  string   `yaml:"required_runtime_class"`
	RequirePolicyDecision bool     `yaml:"require_policy_decision"`
	RequireProvenance     bool     `yaml:"require_provenance"`
	RequireTaintScan      bool     `yaml:"require_taint_scan"`
	MaxResponseBytes      int64    `yaml:"max_response_bytes"`
	TimeoutMS             int64    `yaml:"timeout_ms"`
	RedirectsAllowed      bool     `yaml:"redirects_allowed"`
}

// SourceContracts contains connector projection contracts.
type SourceContracts struct {
	Version    int              `yaml:"version"`
	PolicyMode string           `yaml:"policy_mode"`
	Contracts  []SourceContract `yaml:"contracts"`
}

// SourceContract is a source-specific input/output allowlist.
type SourceContract struct {
	ContractID               string   `yaml:"contract_id"`
	ConnectorID              string   `yaml:"connector_id"`
	Version                  string   `yaml:"version"`
	InputFields              []string `yaml:"input_fields"`
	OutputFields             []string `yaml:"output_fields"`
	ForbiddenOutputFields    []string `yaml:"forbidden_output_fields"`
	RequiredProvenanceFields []string `yaml:"required_provenance_fields"`
	MaxRecords               int      `yaml:"max_records"`
	MaxStringChars           int      `yaml:"max_string_chars"`
}

// Config contains validated embedded Layer 7 configuration.
type Config struct {
	Registry        Registry
	SourceContracts SourceContracts
	MCPManifest     *schemavalidation.Validator
}

// Load validates strict YAML, semantic policy, and the MCP manifest schema.
func Load() (Config, error) {
	var config Config
	if err := decodeStrict(connectorRegistryYAML, &config.Registry); err != nil {
		return Config{}, fmt.Errorf("decode connector registry: %w", err)
	}
	if err := decodeStrict(sourceContractsYAML, &config.SourceContracts); err != nil {
		return Config{}, fmt.Errorf("decode source contracts: %w", err)
	}
	validator, err := schemavalidation.Compile(mcpServerManifestSchemaJSON)
	if err != nil {
		return Config{}, fmt.Errorf("compile mcp manifest schema: %w", err)
	}
	config.MCPManifest = validator
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

// Validate enforces connector and contract security invariants.
func (c Config) Validate() error {
	if c.Registry.Version != 2 ||
		c.Registry.PolicyMode != "fail_closed" ||
		c.Registry.DefaultStatus != "blocked" ||
		len(c.Registry.Connectors) == 0 ||
		c.SourceContracts.Version != 2 ||
		c.SourceContracts.PolicyMode != "fail_closed" ||
		len(c.SourceContracts.Contracts) == 0 ||
		c.MCPManifest == nil {
		return errors.New("layer 7 config is incomplete")
	}
	connectors := make(map[string]Connector, len(c.Registry.Connectors))
	for _, connector := range c.Registry.Connectors {
		if connector.ConnectorID == "" ||
			connector.Status != "active" ||
			connector.OwnerTeam == "" ||
			!connector.ReadOnly ||
			len(connector.AllowedActionClasses) == 0 ||
			connector.RequiredRuntimeClass == "" ||
			!connector.RequirePolicyDecision ||
			!connector.RequireProvenance ||
			!connector.RequireTaintScan ||
			connector.MaxResponseBytes < 1 ||
			connector.MaxResponseBytes > 32<<20 ||
			connector.TimeoutMS < 1000 ||
			connector.TimeoutMS > 120000 ||
			connector.RedirectsAllowed {
			return errors.New("connector registry contains unsafe connector")
		}
		if _, duplicate := connectors[connector.ConnectorID]; duplicate {
			return errors.New("connector registry contains duplicate connector")
		}
		connectors[connector.ConnectorID] = connector
	}
	for _, contract := range c.SourceContracts.Contracts {
		if contract.ContractID == "" ||
			contract.Version == "" ||
			len(contract.InputFields) == 0 ||
			len(contract.OutputFields) == 0 ||
			len(contract.ForbiddenOutputFields) == 0 ||
			len(contract.RequiredProvenanceFields) == 0 ||
			contract.MaxRecords < 1 ||
			contract.MaxStringChars < 1 {
			return errors.New("source contract is incomplete")
		}
		if _, found := connectors[contract.ConnectorID]; !found {
			return errors.New("source contract references unknown connector")
		}
	}
	return nil
}

func decodeStrict(data []byte, destination any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("configuration contains multiple yaml documents")
	}
	return nil
}
