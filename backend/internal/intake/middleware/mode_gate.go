package middleware

import (
	"errors"

	"aeolyzer/internal/intake/contracts"
)

var (
	// ErrModeNotAllowed indicates the derived intent is fundamentally incompatible with the active orchestration mode.
	ErrModeNotAllowed = errors.New("MODE_NOT_ALLOWED")
	// ErrWriteModeRequired blocks state mutations when operating in read-only/plan modes.
	ErrWriteModeRequired = errors.New("WRITE_MODE_REQUIRED")
	// ErrEditModeRequired blocks localized patching unless strictly in Edit mode.
	ErrEditModeRequired = errors.New("EDIT_MODE_REQUIRED")
)

// Resolves the execution envelope for the current session.
// Uses user-supplied mode flag if present; otherwise, infers deterministic modes based on the intent.
// Note: Write/Edit modes strictly require upstream policy approval and are not passively inferred here.
func DeriveMode(input contracts.SanitizedInput, intent contracts.Intent, metadata map[string]interface{}) (contracts.OrchestrationMode, error) {
	if modeStr, ok := metadata["mode"].(string); ok {
		return contracts.OrchestrationMode(modeStr), nil
	}

	switch intent {
	case contracts.IntentDraftArticle:
		// Do not infer write mode automatically; must be explicitly passed in metadata.
		// Protects against accidental execution of generative writes during a planning session.
		return "", ErrWriteModeRequired
	case contracts.IntentEditExisting:
		// Similar to write mode; requires explicit invocation to prevent arbitrary AST updates.
		return "", ErrEditModeRequired
	case contracts.IntentArticlePlanning, contracts.IntentTopicDiscovery:
		return contracts.ModePlan, nil
	case contracts.IntentOptimizeContent:
		return contracts.ModeOptimize, nil
	default:
		return contracts.ModePlan, nil
	}
}

// Verifies that a derived intent is authorized to execute within the current mode.
// Acts as a secondary validation pass before handing off the DAG to Layer 3.
func ValidateModeForIntent(intent contracts.Intent, mode contracts.OrchestrationMode) error {
	switch intent {
	case contracts.IntentDraftArticle:
		if mode != contracts.ModeWrite {
			return ErrWriteModeRequired
		}
	case contracts.IntentEditExisting:
		if mode != contracts.ModeEdit {
			return ErrEditModeRequired
		}
	case contracts.IntentArticlePlanning, contracts.IntentTopicDiscovery:
		// Prevents an agent from hallucinating generative writes while ostensibly only planning.
		if mode != contracts.ModePlan {
			return ErrModeNotAllowed
		}
	case contracts.IntentOptimizeContent:
		if mode != contracts.ModeOptimize && mode != contracts.ModeEdit {
			return ErrModeNotAllowed
		}
	}
	return nil
}

func RequiresWriteMode(intent contracts.Intent) bool {
	return intent == contracts.IntentDraftArticle
}

func RequiresEditMode(intent contracts.Intent) bool {
	return intent == contracts.IntentEditExisting
}
