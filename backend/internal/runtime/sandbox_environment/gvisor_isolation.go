// Package sandboxenvironment verifies kernel-level runtime isolation.
package sandboxenvironment

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var imageDigestPattern = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

// SandboxSpec is a bounded request for an isolated runtime.
type SandboxSpec struct {
	TraceID        string
	TaskID         string
	RuntimeClass   string
	ImageDigest    string
	WallClockLimit time.Duration
	MemoryMB       int64
	MaxProcesses   int64
	NetworkMode    string
}

// DriverHandle is the opaque runtime handle returned by a sandbox driver.
type DriverHandle struct {
	ID string
}

// Attestation is independently queried runtime security evidence.
type Attestation struct {
	RuntimeName        string
	KernelIsolation    bool
	Rootless           bool
	SeccompEnforced    bool
	NetworkNamespace   bool
	CgroupLimits       bool
	AmbientCredentials bool
	ImageDigest        string
}

// Lease is a verified sandbox lease. It contains no credentials.
type Lease struct {
	ID           string
	TraceID      string
	TaskID       string
	RuntimeClass string
	ImageDigest  string
	ExpiresAt    time.Time
}

// Driver creates, attests, and destroys one sandbox.
type Driver interface {
	Create(context.Context, SandboxSpec) (DriverHandle, error)
	Attest(context.Context, string) (Attestation, error)
	Destroy(context.Context, string) error
}

// IsolationController admits only attested gVisor sandboxes.
type IsolationController struct {
	driver Driver
	now    func() time.Time
}

// NewIsolationController creates a fail-closed controller.
func NewIsolationController(
	driver Driver,
	now func() time.Time,
) (*IsolationController, error) {
	if driver == nil || now == nil {
		return nil, errors.New("sandbox isolation controller is not configured")
	}
	return &IsolationController{driver: driver, now: now}, nil
}

// Acquire creates a sandbox and destroys it immediately if attestation fails.
func (c *IsolationController) Acquire(
	ctx context.Context,
	spec SandboxSpec,
) (Lease, error) {
	if c == nil || c.driver == nil || c.now == nil {
		return Lease{}, errors.New("sandbox isolation controller is not configured")
	}
	if err := validateSandboxSpec(spec); err != nil {
		return Lease{}, err
	}
	if err := ctx.Err(); err != nil {
		return Lease{}, fmt.Errorf("acquire sandbox: %w", err)
	}
	handle, err := c.driver.Create(ctx, spec)
	if err != nil {
		return Lease{}, fmt.Errorf("create sandbox: %w", err)
	}
	if handle.ID == "" {
		return Lease{}, errors.New("sandbox driver returned an empty handle")
	}
	attestation, err := c.driver.Attest(ctx, handle.ID)
	if err != nil {
		_ = c.driver.Destroy(ctx, handle.ID)
		return Lease{}, fmt.Errorf("attest sandbox: %w", err)
	}
	if err := validateAttestation(attestation, spec); err != nil {
		_ = c.driver.Destroy(ctx, handle.ID)
		return Lease{}, err
	}
	return Lease{
		ID:           handle.ID,
		TraceID:      spec.TraceID,
		TaskID:       spec.TaskID,
		RuntimeClass: spec.RuntimeClass,
		ImageDigest:  spec.ImageDigest,
		ExpiresAt:    c.now().Add(spec.WallClockLimit),
	}, nil
}

// Release destroys a sandbox lease.
func (c *IsolationController) Release(ctx context.Context, lease Lease) error {
	if c == nil || c.driver == nil {
		return errors.New("sandbox isolation controller is not configured")
	}
	if lease.ID == "" {
		return errors.New("sandbox lease id is required")
	}
	if err := c.driver.Destroy(ctx, lease.ID); err != nil {
		return fmt.Errorf("destroy sandbox: %w", err)
	}
	return nil
}

func validateSandboxSpec(spec SandboxSpec) error {
	if spec.TraceID == "" ||
		spec.TaskID == "" ||
		spec.RuntimeClass == "" ||
		!imageDigestPattern.MatchString(spec.ImageDigest) ||
		spec.WallClockLimit < time.Second ||
		spec.WallClockLimit > 15*time.Minute ||
		spec.MemoryMB < 64 ||
		spec.MemoryMB > 8192 ||
		spec.MaxProcesses < 1 ||
		spec.MaxProcesses > 256 {
		return errors.New("sandbox specification is invalid")
	}
	switch spec.NetworkMode {
	case "none", "egress_proxy_only", "connector_proxy_only":
	default:
		return errors.New("sandbox network mode is invalid")
	}
	return nil
}

func validateAttestation(attestation Attestation, spec SandboxSpec) error {
	if attestation.RuntimeName != "gvisor" ||
		!attestation.KernelIsolation ||
		!attestation.Rootless ||
		!attestation.SeccompEnforced ||
		!attestation.NetworkNamespace ||
		!attestation.CgroupLimits ||
		attestation.AmbientCredentials ||
		attestation.ImageDigest != spec.ImageDigest {
		return errors.New("sandbox attestation failed")
	}
	return nil
}
