package extensions

import (
	_ "embed"
	"errors"
	"fmt"

	"aeolyzer/internal/schemavalidation"
)

var (
	//go:embed presentation.schema.json
	presentationSchemaJSON []byte
	//go:embed a2ui-frame.schema.json
	a2uiFrameSchemaJSON []byte
	//go:embed a2ui-catalog.schema.json
	a2uiCatalogSchemaJSON []byte
	//go:embed ui-event.schema.json
	uiEventSchemaJSON []byte
	//go:embed approval.schema.json
	approvalSchemaJSON []byte
	//go:embed surface-patch.schema.json
	surfacePatchSchemaJSON []byte
	//go:embed a2a-agent-card.schema.json
	a2aAgentCardSchemaJSON []byte
	//go:embed a2a-envelope.schema.json
	a2aEnvelopeSchemaJSON []byte
	//go:embed catalog-lock.yaml
	catalogLockYAML []byte
)

// Contract identifies one Layer 5 schema.
type Contract string

const (
	ContractPresentation Contract = "presentation"
	ContractA2UIFrame    Contract = "a2ui_frame"
	ContractA2UICatalog  Contract = "a2ui_catalog"
	ContractUIEvent      Contract = "ui_event"
	ContractApproval     Contract = "approval"
	ContractSurfacePatch Contract = "surface_patch"
	ContractA2AAgentCard Contract = "a2a_agent_card"
	ContractA2AEnvelope  Contract = "a2a_envelope"
)

// Schemas contains compiled Layer 5 contracts.
type Schemas struct {
	validators map[Contract]*schemavalidation.Validator
}

// NewSchemas compiles all Layer 5 JSON Schemas and validates the embedded
// catalog lock.
func NewSchemas() (*Schemas, error) {
	inputs := map[Contract][]byte{
		ContractPresentation: presentationSchemaJSON,
		ContractA2UIFrame:    a2uiFrameSchemaJSON,
		ContractA2UICatalog:  a2uiCatalogSchemaJSON,
		ContractUIEvent:      uiEventSchemaJSON,
		ContractApproval:     approvalSchemaJSON,
		ContractSurfacePatch: surfacePatchSchemaJSON,
		ContractA2AAgentCard: a2aAgentCardSchemaJSON,
		ContractA2AEnvelope:  a2aEnvelopeSchemaJSON,
	}
	validators := make(map[Contract]*schemavalidation.Validator, len(inputs))
	for contract, schema := range inputs {
		validator, err := schemavalidation.Compile(schema)
		if err != nil {
			return nil, fmt.Errorf("compile %s schema: %w", contract, err)
		}
		validators[contract] = validator
	}
	schemas := &Schemas{validators: validators}
	if err := schemas.ValidateYAML(ContractA2UICatalog, catalogLockYAML); err != nil {
		return nil, fmt.Errorf("validate embedded catalog lock: %w", err)
	}
	return schemas, nil
}

// ValidateJSON validates one Layer 5 JSON contract.
func (s *Schemas) ValidateJSON(contract Contract, data []byte) error {
	if s == nil {
		return errors.New("layer 5 schemas are not configured")
	}
	validator := s.validators[contract]
	if validator == nil {
		return errors.New("layer 5 contract is unknown")
	}
	return validator.ValidateJSON(data)
}

// ValidateYAML validates one YAML document against a Layer 5 JSON Schema.
func (s *Schemas) ValidateYAML(contract Contract, data []byte) error {
	if s == nil {
		return errors.New("layer 5 schemas are not configured")
	}
	validator := s.validators[contract]
	if validator == nil {
		return errors.New("layer 5 contract is unknown")
	}
	return validator.ValidateYAML(data)
}

// CatalogLock returns a defensive copy of the validated embedded catalog.
func CatalogLock() []byte {
	return append([]byte(nil), catalogLockYAML...)
}
