// Package supplychaindefense enforces pinned, pre-approved dependencies.
package supplychaindefense

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// ErrDynamicInstallBlocked indicates an attempted runtime dependency change.
	ErrDynamicInstallBlocked = errors.New("dynamic dependency installation is blocked")
	// ErrShellInvocationBlocked indicates an unstructured shell boundary.
	ErrShellInvocationBlocked = errors.New("shell command execution is blocked")
	// ErrDependencyNotPinned indicates a dependency missing an exact allowlist
	// version and digest.
	ErrDependencyNotPinned = errors.New("dependency is not pinned and approved")
)

var installTextPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(?:^|[\s;&|])(?:npm|pnpm|yarn|bun)\s+(?:add|install|i)(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])(?:pip|pip3|uv)\s+install(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])python(?:3)?\s+-m\s+pip\s+install(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])go\s+(?:get|install)(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])cargo\s+(?:add|install)(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])gem\s+install(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:^|[\s;&|])(?:apt|apt-get|apk|dnf|yum|pacman|brew|choco|winget)\s+(?:add|install)(?:\s|$)`),
	regexp.MustCompile(`(?i)(?:curl|wget).{0,200}\|\s*(?:sh|bash|zsh|powershell)`),
}

// Dependency identifies one locked dependency.
type Dependency struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Digest    string `json:"digest"`
}

// Policy contains exact approved dependencies.
type Policy struct {
	Dependencies []Dependency `json:"dependencies"`
}

// Guard blocks dynamic installation and validates pinned dependencies.
type Guard struct {
	approved map[string]Dependency
}

// NewGuard validates and indexes an exact dependency allowlist.
func NewGuard(policy Policy) (*Guard, error) {
	approved := make(map[string]Dependency, len(policy.Dependencies))
	for _, dependency := range policy.Dependencies {
		if err := validateDependency(dependency); err != nil {
			return nil, err
		}
		key := dependencyKey(dependency.Ecosystem, dependency.Name)
		if _, duplicate := approved[key]; duplicate {
			return nil, fmt.Errorf("duplicate dependency %q", key)
		}
		approved[key] = dependency
	}
	return &Guard{approved: approved}, nil
}

// ValidateInvocation rejects shells, package managers, and installer
// subcommands. Callers must pass a structured executable and argument vector.
func (g *Guard) ValidateInvocation(program string, args []string) error {
	if g == nil {
		return errors.New("dependency guard is not configured")
	}
	executable := strings.ToLower(strings.TrimSuffix(filepath.Base(program), ".exe"))
	if executable == "" {
		return errors.New("executable is required")
	}
	if isShell(executable) {
		return ErrShellInvocationBlocked
	}
	if isPackageManagerInstall(executable, args) {
		return ErrDynamicInstallBlocked
	}
	return nil
}

// ScanGeneratedText rejects installer sequences in generated scripts before
// they enter Layer 6 execution.
func (g *Guard) ScanGeneratedText(text string) error {
	if g == nil {
		return errors.New("dependency guard is not configured")
	}
	if len(text) > 1<<20 {
		return errors.New("generated script exceeds scan limit")
	}
	for _, pattern := range installTextPatterns {
		if pattern.MatchString(text) {
			return ErrDynamicInstallBlocked
		}
	}
	return nil
}

// AuthorizeDependency requires an exact ecosystem, name, version, and digest
// match.
func (g *Guard) AuthorizeDependency(dependency Dependency) error {
	if g == nil {
		return errors.New("dependency guard is not configured")
	}
	approved, found := g.approved[dependencyKey(dependency.Ecosystem, dependency.Name)]
	if !found ||
		dependency.Version != approved.Version ||
		!constantTimeStringEqual(dependency.Digest, approved.Digest) {
		return ErrDependencyNotPinned
	}
	return nil
}

func validateDependency(dependency Dependency) error {
	if dependency.Ecosystem == "" ||
		dependency.Name == "" ||
		dependency.Version == "" ||
		!strings.HasPrefix(dependency.Digest, "sha256:") ||
		len(dependency.Digest) != len("sha256:")+64 {
		return ErrDependencyNotPinned
	}
	if strings.ContainsAny(dependency.Name, " \t\r\n/\\") {
		return ErrDependencyNotPinned
	}
	return nil
}

func isShell(executable string) bool {
	switch executable {
	case "sh", "bash", "zsh", "fish", "cmd", "powershell", "pwsh":
		return true
	default:
		return false
	}
}

func isPackageManagerInstall(executable string, args []string) bool {
	normalized := make([]string, len(args))
	for i, arg := range args {
		normalized[i] = strings.ToLower(strings.TrimSpace(arg))
	}
	first := ""
	if len(normalized) > 0 {
		first = normalized[0]
	}
	switch executable {
	case "npm", "pnpm", "yarn", "bun":
		return first == "add" || first == "install" || first == "i"
	case "pip", "pip3", "uv":
		return first == "install" || containsSequence(normalized, "pip", "install")
	case "python", "python3":
		return containsSequence(normalized, "-m", "pip", "install")
	case "go":
		return first == "get" || first == "install"
	case "cargo":
		return first == "add" || first == "install"
	case "gem":
		return first == "install"
	case "apt", "apt-get", "apk", "dnf", "yum", "pacman", "brew", "choco", "winget":
		return first == "add" || first == "install"
	default:
		return false
	}
}

func containsSequence(values []string, sequence ...string) bool {
	if len(sequence) == 0 || len(values) < len(sequence) {
		return false
	}
	for start := 0; start <= len(values)-len(sequence); start++ {
		match := true
		for offset, expected := range sequence {
			if values[start+offset] != expected {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func dependencyKey(ecosystem, name string) string {
	return strings.ToLower(ecosystem) + ":" + strings.ToLower(name)
}

func constantTimeStringEqual(left, right string) bool {
	if len(left) != len(right) {
		return false
	}
	var difference byte
	for i := range left {
		difference |= left[i] ^ right[i]
	}
	return difference == 0
}
