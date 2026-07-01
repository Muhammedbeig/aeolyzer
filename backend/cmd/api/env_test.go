package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalEnvironmentFromExplicitPath(t *testing.T) {
	const name = "AEOLYZER_ENV_TEST_LOCAL_FILE"
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(name+"=loaded\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AEOLYZER_ENV_FILE", path)
	previous, existed := os.LookupEnv(name)
	if err := os.Unsetenv(name); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(name, previous)
			return
		}
		_ = os.Unsetenv(name)
	})

	if err := loadLocalEnvironment(); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv(name); got != "loaded" {
		t.Fatalf("%s = %q, want loaded", name, got)
	}
}
