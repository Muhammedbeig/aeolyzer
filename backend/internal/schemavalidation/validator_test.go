package schemavalidation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidatorAcceptsValidJSONAndYAML(t *testing.T) {
	validator := compileTestValidator(t)

	for name, test := range map[string]struct {
		validate func([]byte) error
		input    string
	}{
		"json": {
			validate: validator.ValidateJSON,
			input:    `{"version":2,"enabled":true}`,
		},
		"yaml": {
			validate: validator.ValidateYAML,
			input:    "version: 2\nenabled: true\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := test.validate([]byte(test.input)); err != nil {
				t.Fatalf("validate valid %s document: %v", name, err)
			}
		})
	}
}

func TestValidatorRejectsMalformedAndInvalidDocuments(t *testing.T) {
	validator := compileTestValidator(t)

	for name, test := range map[string]struct {
		validate func([]byte) error
		input    string
	}{
		"malformed json": {
			validate: validator.ValidateJSON,
			input:    `{"version":`,
		},
		"invalid json": {
			validate: validator.ValidateJSON,
			input:    `{"version":1,"enabled":true}`,
		},
		"malformed yaml": {
			validate: validator.ValidateYAML,
			input:    "version: [\n",
		},
		"invalid yaml": {
			validate: validator.ValidateYAML,
			input:    "version: 2\nenabled: yes\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := test.validate([]byte(test.input)); err == nil {
				t.Fatalf("validate %s document returned nil error", name)
			}
		})
	}
}

func TestCompileRejectsEmptyInvalidAndOversizedSchemas(t *testing.T) {
	tests := map[string]struct {
		input []byte
		want  error
	}{
		"empty": {
			input: nil,
		},
		"invalid": {
			input: []byte(`{"type":`),
		},
		"oversized": {
			input: []byte(strings.Repeat("x", MaxSchemaBytes+1)),
			want:  ErrSchemaTooLarge,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Compile(test.input)
			if err == nil {
				t.Fatal("Compile() returned nil error")
			}
			if test.want != nil && !errors.Is(err, test.want) {
				t.Fatalf("Compile() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestValidatorRejectsOversizedDocument(t *testing.T) {
	validator := compileTestValidator(t)
	document := []byte(strings.Repeat(" ", MaxDocumentBytes+1))

	if err := validator.ValidateJSON(document); !errors.Is(err, ErrDocumentTooLarge) {
		t.Fatalf("ValidateJSON() error = %v, want %v", err, ErrDocumentTooLarge)
	}
}

func compileTestValidator(t *testing.T) *Validator {
	t.Helper()
	validator, err := Compile([]byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"additionalProperties": false,
		"required": ["version", "enabled"],
		"properties": {
			"version": {"const": 2},
			"enabled": {"type": "boolean"}
		}
	}`))
	if err != nil {
		t.Fatalf("Compile() failed: %v", err)
	}
	return validator
}
