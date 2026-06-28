package a2uitranslator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"aeolyzer/internal/extensions"
	"go.yaml.in/yaml/v3"
)

// Catalog is the runtime component allowlist.
type Catalog struct {
	Version        int                `yaml:"version"`
	PolicyMode     string             `yaml:"policy_mode"`
	CatalogID      string             `yaml:"catalog_id"`
	CatalogVersion string             `yaml:"catalog_version"`
	SchemaVersion  string             `yaml:"schema_version"`
	Components     []CatalogComponent `yaml:"components"`
}

// CatalogComponent defines allowed props and child behavior.
type CatalogComponent struct {
	Type           string   `yaml:"type"`
	AllowedProps   []string `yaml:"allowed_props"`
	RequiredProps  []string `yaml:"required_props"`
	AllowsChildren bool     `yaml:"allows_children"`
	MaxChildren    int      `yaml:"max_children"`
}

// LoadCatalog strictly loads the embedded catalog lock.
func LoadCatalog(data []byte, schemas *extensions.Schemas) (Catalog, error) {
	if schemas == nil {
		return Catalog{}, errors.New("a2ui schemas are required")
	}
	if err := schemas.ValidateYAML(extensions.ContractA2UICatalog, data); err != nil {
		return Catalog{}, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var catalog Catalog
	if err := decoder.Decode(&catalog); err != nil {
		return Catalog{}, fmt.Errorf("decode a2ui catalog: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Catalog{}, errors.New("a2ui catalog contains multiple documents")
	}
	return catalog, nil
}

// ValidateCatalogNodes enforces the locked component and prop allowlist.
func ValidateCatalogNodes(nodes []extensions.A2UINode, catalog Catalog) error {
	components := make(map[string]CatalogComponent, len(catalog.Components))
	for _, component := range catalog.Components {
		if _, duplicate := components[component.Type]; duplicate {
			return errors.New("a2ui catalog contains duplicate component")
		}
		components[component.Type] = component
	}
	for _, node := range nodes {
		component, found := components[node.Type]
		if !found {
			return fmt.Errorf("a2ui component %q is not in the catalog", node.Type)
		}
		if !component.AllowsChildren && len(node.Children) > 0 {
			return fmt.Errorf("a2ui component %q does not allow children", node.Type)
		}
		if len(node.Children) > component.MaxChildren {
			return fmt.Errorf("a2ui component %q exceeds child limit", node.Type)
		}
		allowed := stringSet(component.AllowedProps)
		for key, value := range node.Props {
			if _, ok := allowed[key]; !ok {
				return fmt.Errorf("a2ui prop %q is not allowed for %q", key, node.Type)
			}
			if executableValue(key, value) {
				return fmt.Errorf("a2ui prop %q contains executable content", key)
			}
		}
		for _, required := range component.RequiredProps {
			if _, found := node.Props[required]; !found {
				return fmt.Errorf("a2ui component %q is missing prop %q", node.Type, required)
			}
		}
	}
	return nil
}

func executableValue(key string, value any) bool {
	lowerKey := strings.ToLower(key)
	if strings.HasPrefix(lowerKey, "on") ||
		strings.Contains(lowerKey, "script") ||
		strings.Contains(lowerKey, "html") ||
		strings.Contains(lowerKey, "style") ||
		strings.Contains(lowerKey, "class") {
		return true
	}
	text, ok := value.(string)
	if !ok {
		return false
	}
	lower := strings.ToLower(text)
	return strings.Contains(lower, "<script") ||
		strings.Contains(lower, "javascript:") ||
		strings.Contains(lower, "data:text/html") ||
		strings.Contains(lower, "expression(")
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
