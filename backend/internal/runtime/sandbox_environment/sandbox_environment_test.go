package sandboxenvironment

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type driverStub struct {
	attestation Attestation
	destroyed   []string
}

func (d *driverStub) Create(_ context.Context, _ SandboxSpec) (DriverHandle, error) {
	return DriverHandle{ID: "sandbox-1"}, nil
}

func (d *driverStub) Attest(_ context.Context, _ string) (Attestation, error) {
	return d.attestation, nil
}

func (d *driverStub) Destroy(_ context.Context, id string) error {
	d.destroyed = append(d.destroyed, id)
	return nil
}

func TestIsolationControllerRequiresCompleteAttestation(t *testing.T) {
	spec := validSandboxSpec()
	driver := &driverStub{attestation: validAttestation(spec.ImageDigest)}
	controller, err := NewIsolationController(driver, time.Now)
	if err != nil {
		t.Fatalf("NewIsolationController() failed: %v", err)
	}
	lease, err := controller.Acquire(context.Background(), spec)
	if err != nil {
		t.Fatalf("IsolationController.Acquire() failed: %v", err)
	}
	if lease.ID != "sandbox-1" {
		t.Fatalf("IsolationController.Acquire().ID = %q, want sandbox-1", lease.ID)
	}
}

func TestIsolationControllerDestroysFailedSandbox(t *testing.T) {
	spec := validSandboxSpec()
	attestation := validAttestation(spec.ImageDigest)
	attestation.Rootless = false
	driver := &driverStub{attestation: attestation}
	controller, err := NewIsolationController(driver, time.Now)
	if err != nil {
		t.Fatalf("NewIsolationController() failed: %v", err)
	}
	if _, err := controller.Acquire(context.Background(), spec); err == nil {
		t.Fatal("IsolationController.Acquire() accepted unsafe attestation")
	}
	if len(driver.destroyed) != 1 || driver.destroyed[0] != "sandbox-1" {
		t.Fatalf("failed sandbox destruction = %v, want sandbox-1", driver.destroyed)
	}
}

func TestStateManagerCreatesAndResetsMarkedWorkspace(t *testing.T) {
	root := t.TempDir()
	manager, err := NewStateManager(root)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}
	path, err := manager.Create("trace-1", "task-1")
	if err != nil {
		t.Fatalf("StateManager.Create() failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, workspaceMarker)); err != nil {
		t.Fatalf("workspace marker missing: %v", err)
	}
	if err := manager.Reset(path); err != nil {
		t.Fatalf("StateManager.Reset() failed: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("workspace still exists: %v", err)
	}
}

func TestStateManagerRejectsRootAndUnmarkedDirectory(t *testing.T) {
	root := t.TempDir()
	manager, err := NewStateManager(root)
	if err != nil {
		t.Fatalf("NewStateManager() failed: %v", err)
	}
	if err := manager.Reset(root); err == nil {
		t.Fatal("StateManager.Reset(root) returned nil error")
	}
	unmarked := filepath.Join(root, "unmarked")
	if err := os.Mkdir(unmarked, 0o700); err != nil {
		t.Fatalf("os.Mkdir() failed: %v", err)
	}
	if err := manager.Reset(unmarked); err == nil {
		t.Fatal("StateManager.Reset(unmarked) returned nil error")
	}
}

func validSandboxSpec() SandboxSpec {
	return SandboxSpec{
		TraceID:        "trace-1",
		TaskID:         "task-1",
		RuntimeClass:   "deterministic_skill_script",
		ImageDigest:    "sha256:" + strings.Repeat("a", 64),
		WallClockLimit: time.Minute,
		MemoryMB:       512,
		MaxProcesses:   16,
		NetworkMode:    "none",
	}
}

func validAttestation(imageDigest string) Attestation {
	return Attestation{
		RuntimeName:        "gvisor",
		KernelIsolation:    true,
		Rootless:           true,
		SeccompEnforced:    true,
		NetworkNamespace:   true,
		CgroupLimits:       true,
		AmbientCredentials: false,
		ImageDigest:        imageDigest,
	}
}
