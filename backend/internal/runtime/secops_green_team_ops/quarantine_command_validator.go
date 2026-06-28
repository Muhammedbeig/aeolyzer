package secops_green_team_ops

import "errors"
import "aeolyzer/internal/runtime"

var ErrInvalidQuarantineCommand = errors.New("INVALID_QUARANTINE_COMMAND")

// ValidateQuarantineCommand ensures that stateful freezes are strictly authorized.
// Layer 6 does not decide to quarantine; it only executes signed commands from Layer 8.
func ValidateQuarantineCommand(cmd runtime.QuarantineCommand) error {
	if cmd.Signature == "" || cmd.TargetScope == "" {
		return ErrInvalidQuarantineCommand
	}
	return nil
}
