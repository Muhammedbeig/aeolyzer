// Package a2uitranslator validates and converts declarative A2UI payloads.
package a2uitranslator

import (
	"encoding/json"
	"errors"
	"fmt"

	"aeolyzer/internal/extensions"
)

// SchemaManager applies JSON Schema and graph-integrity validation.
type SchemaManager struct {
	schemas *extensions.Schemas
}

// NewSchemaManager creates a manager from compiled Layer 5 schemas.
func NewSchemaManager(schemas *extensions.Schemas) (*SchemaManager, error) {
	if schemas == nil {
		return nil, errors.New("a2ui schemas are required")
	}
	return &SchemaManager{schemas: schemas}, nil
}

// ValidateFrame validates schema, unique IDs, references, root, and cycles.
func (m *SchemaManager) ValidateFrame(frame extensions.A2UIFrame) error {
	if m == nil || m.schemas == nil {
		return errors.New("a2ui schema manager is not configured")
	}
	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("encode a2ui frame: %w", err)
	}
	if err := m.schemas.ValidateJSON(extensions.ContractA2UIFrame, data); err != nil {
		return err
	}
	nodes := make(map[string]extensions.A2UINode, len(frame.Nodes))
	for _, node := range frame.Nodes {
		if _, duplicate := nodes[node.ID]; duplicate {
			return fmt.Errorf("duplicate a2ui node id %q", node.ID)
		}
		nodes[node.ID] = node
	}
	if _, found := nodes[frame.RootID]; !found {
		return errors.New("a2ui root node is missing")
	}
	for _, node := range frame.Nodes {
		for _, child := range node.Children {
			if _, found := nodes[child]; !found {
				return fmt.Errorf("a2ui child node %q is missing", child)
			}
		}
	}
	visiting := make(map[string]bool, len(nodes))
	visited := make(map[string]bool, len(nodes))
	var visit func(string, int) error
	visit = func(id string, depth int) error {
		if depth > 64 {
			return errors.New("a2ui graph exceeds depth limit")
		}
		if visiting[id] {
			return errors.New("a2ui graph contains a cycle")
		}
		if visited[id] {
			return nil
		}
		visiting[id] = true
		for _, child := range nodes[id].Children {
			if err := visit(child, depth+1); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		return nil
	}
	if err := visit(frame.RootID, 0); err != nil {
		return err
	}
	if len(visited) != len(nodes) {
		return errors.New("a2ui frame contains unreachable nodes")
	}
	return nil
}
