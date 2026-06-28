package sandboxenvironment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const workspaceMarker = ".aeolyzer-ephemeral-workspace"

// StateManager creates and destroys isolated ephemeral workspaces beneath one
// configured root.
type StateManager struct {
	root string
}

// NewStateManager validates and canonicalizes an existing workspace root.
func NewStateManager(root string) (*StateManager, error) {
	if root == "" {
		return nil, errors.New("ephemeral state root is required")
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve ephemeral state root: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, fmt.Errorf("evaluate ephemeral state root: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return nil, errors.New("ephemeral state root must be an existing directory")
	}
	return &StateManager{root: filepath.Clean(resolved)}, nil
}

// Create creates a mode-0700 workspace with an ownership marker.
func (m *StateManager) Create(traceID, taskID string) (string, error) {
	if m == nil || m.root == "" {
		return "", errors.New("ephemeral state manager is not configured")
	}
	if !safePathToken(traceID) || !safePathToken(taskID) {
		return "", errors.New("ephemeral trace and task ids are invalid")
	}
	path, err := os.MkdirTemp(m.root, traceID+"-"+taskID+"-")
	if err != nil {
		return "", fmt.Errorf("create ephemeral workspace: %w", err)
	}
	if err := os.Chmod(path, 0o700); err != nil {
		_ = os.RemoveAll(path)
		return "", fmt.Errorf("secure ephemeral workspace: %w", err)
	}
	marker := filepath.Join(path, workspaceMarker)
	if err := os.WriteFile(marker, []byte(traceID+"\n"+taskID+"\n"), 0o600); err != nil {
		_ = os.RemoveAll(path)
		return "", fmt.Errorf("mark ephemeral workspace: %w", err)
	}
	return path, nil
}

// Reset removes a marked workspace after proving it remains under the
// configured root and is not a symlink.
func (m *StateManager) Reset(path string) error {
	if m == nil || m.root == "" {
		return errors.New("ephemeral state manager is not configured")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve ephemeral workspace: %w", err)
	}
	clean := filepath.Clean(absolute)
	if clean == m.root || !isWithin(m.root, clean) {
		return errors.New("ephemeral workspace escapes configured root")
	}
	info, err := os.Lstat(clean)
	if err != nil {
		return fmt.Errorf("inspect ephemeral workspace: %w", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("ephemeral workspace is not a real directory")
	}
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		return fmt.Errorf("evaluate ephemeral workspace: %w", err)
	}
	if resolved != clean || !isWithin(m.root, resolved) {
		return errors.New("ephemeral workspace contains an unsafe path binding")
	}
	marker := filepath.Join(clean, workspaceMarker)
	markerInfo, err := os.Lstat(marker)
	if err != nil || !markerInfo.Mode().IsRegular() {
		return errors.New("ephemeral workspace ownership marker is missing")
	}
	if err := os.RemoveAll(clean); err != nil {
		return fmt.Errorf("remove ephemeral workspace: %w", err)
	}
	return nil
}

func safePathToken(value string) bool {
	if value == "" || len(value) > 64 {
		return false
	}
	for _, character := range value {
		if (character < 'a' || character > 'z') &&
			(character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') &&
			character != '-' &&
			character != '_' {
			return false
		}
	}
	return true
}

func isWithin(root, target string) bool {
	prefix := root + string(filepath.Separator)
	return strings.HasPrefix(
		strings.ToLower(target),
		strings.ToLower(prefix),
	)
}
