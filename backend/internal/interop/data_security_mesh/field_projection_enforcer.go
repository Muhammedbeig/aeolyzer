package datasecuritymesh

import (
	"errors"
	"fmt"
	"strings"
)

// ProjectionPolicy is a connector-specific top-level field allowlist.
type ProjectionPolicy struct {
	AllowedFields  []string
	RequiredFields []string
	MaxFields      int
}

// EnforceProjection returns a defensive copy containing only explicitly
// authorized fields.
func EnforceProjection(
	source map[string]any,
	requested []string,
	policy ProjectionPolicy,
) (map[string]any, error) {
	if source == nil ||
		len(requested) == 0 ||
		len(policy.AllowedFields) == 0 ||
		policy.MaxFields < 1 ||
		len(requested) > policy.MaxFields {
		return nil, errors.New("field projection request is invalid")
	}
	allowed := make(map[string]struct{}, len(policy.AllowedFields))
	for _, field := range policy.AllowedFields {
		if !safeFieldName(field) {
			return nil, errors.New("field projection policy is invalid")
		}
		allowed[field] = struct{}{}
	}
	required := make(map[string]struct{}, len(policy.RequiredFields))
	for _, field := range policy.RequiredFields {
		if _, ok := allowed[field]; !ok {
			return nil, errors.New("required projection field is not allowed")
		}
		required[field] = struct{}{}
	}

	result := make(map[string]any, len(requested))
	seen := make(map[string]struct{}, len(requested))
	for _, field := range requested {
		if !safeFieldName(field) {
			return nil, errors.New("requested projection field is invalid")
		}
		if _, duplicate := seen[field]; duplicate {
			return nil, errors.New("requested projection contains duplicate field")
		}
		seen[field] = struct{}{}
		if _, ok := allowed[field]; !ok {
			return nil, fmt.Errorf("field %q is not authorized for projection", field)
		}
		value, found := source[field]
		if !found {
			return nil, fmt.Errorf("projected field %q is absent", field)
		}
		result[field] = deepCopyValue(value)
	}
	for field := range required {
		if _, included := result[field]; !included {
			return nil, fmt.Errorf("required projection field %q was omitted", field)
		}
	}
	return result, nil
}

func safeFieldName(field string) bool {
	if field == "" || len(field) > 64 || strings.ContainsAny(field, ".[]/\\ \t\r\n") {
		return false
	}
	for _, character := range field {
		if (character < 'a' || character > 'z') &&
			(character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') &&
			character != '_' {
			return false
		}
	}
	return true
}

func deepCopyValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, child := range typed {
			result[key] = deepCopyValue(child)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = deepCopyValue(child)
		}
		return result
	default:
		return typed
	}
}
