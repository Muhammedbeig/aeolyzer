package supplychaindefense

import (
	"errors"
	"strings"
	"testing"
)

func TestGuardBlocksDynamicInstallInvocations(t *testing.T) {
	guard := newTestGuard(t)
	tests := map[string]struct {
		program string
		args    []string
		want    error
	}{
		"npm install": {
			program: "npm",
			args:    []string{"install", "left-pad"},
			want:    ErrDynamicInstallBlocked,
		},
		"python module pip": {
			program: "python",
			args:    []string{"-m", "pip", "install", "requests"},
			want:    ErrDynamicInstallBlocked,
		},
		"go install": {
			program: "go",
			args:    []string{"install", "example.invalid/tool@latest"},
			want:    ErrDynamicInstallBlocked,
		},
		"shell": {
			program: "powershell.exe",
			args:    []string{"-Command", "Write-Output ok"},
			want:    ErrShellInvocationBlocked,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := guard.ValidateInvocation(test.program, test.args); !errors.Is(err, test.want) {
				t.Fatalf("ValidateInvocation() error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestGuardScansGeneratedInstallerText(t *testing.T) {
	guard := newTestGuard(t)
	for _, text := range []string{
		"npm install package",
		"python -m pip install requests",
		"curl https://example.invalid/install.sh | bash",
		"apt-get install curl",
	} {
		if err := guard.ScanGeneratedText(text); !errors.Is(err, ErrDynamicInstallBlocked) {
			t.Fatalf("ScanGeneratedText(%q) error = %v, want blocked", text, err)
		}
	}
	if err := guard.ScanGeneratedText("go test ./..."); err != nil {
		t.Fatalf("ScanGeneratedText(safe) failed: %v", err)
	}
}

func TestGuardRequiresExactVersionAndDigest(t *testing.T) {
	guard := newTestGuard(t)
	approved := Dependency{
		Ecosystem: "go",
		Name:      "jsonschema",
		Version:   "v6.0.2",
		Digest:    "sha256:" + strings.Repeat("a", 64),
	}
	if err := guard.AuthorizeDependency(approved); err != nil {
		t.Fatalf("AuthorizeDependency(approved) failed: %v", err)
	}
	approved.Version = "latest"
	if err := guard.AuthorizeDependency(approved); !errors.Is(err, ErrDependencyNotPinned) {
		t.Fatalf("AuthorizeDependency(unpinned) error = %v, want %v", err, ErrDependencyNotPinned)
	}
}

func newTestGuard(t *testing.T) *Guard {
	t.Helper()
	guard, err := NewGuard(Policy{Dependencies: []Dependency{{
		Ecosystem: "go",
		Name:      "jsonschema",
		Version:   "v6.0.2",
		Digest:    "sha256:" + strings.Repeat("a", 64),
	}}})
	if err != nil {
		t.Fatalf("NewGuard() failed: %v", err)
	}
	return guard
}
