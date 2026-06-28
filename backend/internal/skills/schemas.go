package skills

import (
	_ "embed"
	"fmt"

	"aeolyzer/internal/schemavalidation"
)

var (
	//go:embed registry.schema.json
	registrySchemaJSON []byte
	//go:embed skill.schema.json
	skillSchemaJSON []byte
	//go:embed resource-manifest.schema.json
	resourceManifestSchemaJSON []byte
	//go:embed eval-manifest.schema.json
	evalManifestSchemaJSON []byte
)

// Schemas contains compiled Layer 4 artifact contracts.
type Schemas struct {
	registry         *schemavalidation.Validator
	skill            *schemavalidation.Validator
	resourceManifest *schemavalidation.Validator
	evalManifest     *schemavalidation.Validator
}

// NewSchemas compiles every Layer 4 schema and fails if any contract is invalid.
func NewSchemas() (*Schemas, error) {
	registry, err := schemavalidation.Compile(registrySchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("compile skill registry schema: %w", err)
	}
	skill, err := schemavalidation.Compile(skillSchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("compile skill schema: %w", err)
	}
	resourceManifest, err := schemavalidation.Compile(resourceManifestSchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("compile resource manifest schema: %w", err)
	}
	evalManifest, err := schemavalidation.Compile(evalManifestSchemaJSON)
	if err != nil {
		return nil, fmt.Errorf("compile eval manifest schema: %w", err)
	}
	return &Schemas{
		registry:         registry,
		skill:            skill,
		resourceManifest: resourceManifest,
		evalManifest:     evalManifest,
	}, nil
}

// ValidateRegistry validates one registry YAML document.
func (s *Schemas) ValidateRegistry(data []byte) error {
	if s == nil || s.registry == nil {
		return ErrSchemasNotConfigured
	}
	return s.registry.ValidateYAML(data)
}

// ValidateSkillFrontmatter validates one SKILL.md frontmatter YAML document.
func (s *Schemas) ValidateSkillFrontmatter(data []byte) error {
	if s == nil || s.skill == nil {
		return ErrSchemasNotConfigured
	}
	return s.skill.ValidateYAML(data)
}

// ValidateResourceManifest validates one resource manifest YAML document.
func (s *Schemas) ValidateResourceManifest(data []byte) error {
	if s == nil || s.resourceManifest == nil {
		return ErrSchemasNotConfigured
	}
	return s.resourceManifest.ValidateYAML(data)
}

// ValidateEvalManifest validates one evaluation manifest YAML document.
func (s *Schemas) ValidateEvalManifest(data []byte) error {
	if s == nil || s.evalManifest == nil {
		return ErrSchemasNotConfigured
	}
	return s.evalManifest.ValidateYAML(data)
}
