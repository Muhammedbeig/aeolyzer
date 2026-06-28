// Package schemavalidation compiles and applies repository JSON Schema
// contracts to JSON and YAML documents.
package schemavalidation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"go.yaml.in/yaml/v3"
)

const (
	// MaxSchemaBytes limits a schema before compilation.
	MaxSchemaBytes = 1 << 20
	// MaxDocumentBytes limits a document before parsing and validation.
	MaxDocumentBytes = 4 << 20
)

var (
	// ErrSchemaTooLarge is returned before compiling an oversized schema.
	ErrSchemaTooLarge = errors.New("schema exceeds size limit")
	// ErrDocumentTooLarge is returned before parsing an oversized document.
	ErrDocumentTooLarge = errors.New("document exceeds size limit")
)

// Validator is an immutable compiled JSON Schema validator.
type Validator struct {
	schema *jsonschema.Schema
}

// Compile compiles a JSON Schema using draft 2020-12 and strict format checks.
func Compile(data []byte) (*Validator, error) {
	if len(data) == 0 {
		return nil, errors.New("schema is empty")
	}
	if len(data) > MaxSchemaBytes {
		return nil, ErrSchemaTooLarge
	}

	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode schema: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	compiler.AssertFormat()
	if err := compiler.AddResource("inmemory://schema.json", document); err != nil {
		return nil, fmt.Errorf("add schema resource: %w", err)
	}
	schema, err := compiler.Compile("inmemory://schema.json")
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}
	return &Validator{schema: schema}, nil
}

// ValidateJSON validates one JSON document.
func (v *Validator) ValidateJSON(data []byte) error {
	if v == nil || v.schema == nil {
		return errors.New("validator is not configured")
	}
	if len(data) > MaxDocumentBytes {
		return ErrDocumentTooLarge
	}

	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode json document: %w", err)
	}
	if err := v.schema.Validate(document); err != nil {
		return fmt.Errorf("validate json document: %w", err)
	}
	return nil
}

// ValidateYAML validates one YAML document after converting its data model to
// the JSON-compatible representation required by JSON Schema.
func (v *Validator) ValidateYAML(data []byte) error {
	if v == nil || v.schema == nil {
		return errors.New("validator is not configured")
	}
	if len(data) > MaxDocumentBytes {
		return ErrDocumentTooLarge
	}

	var document any
	if err := yaml.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("decode yaml document: %w", err)
	}
	normalized, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("normalize yaml document: %w", err)
	}
	return v.ValidateJSON(normalized)
}
