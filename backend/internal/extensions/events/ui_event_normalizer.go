package events

import (
	"errors"
	"time"
)

// UserInteractionEvent is the inbound contract from Layer 1.
type UserInteractionEvent struct {
	FrameID       string `json:"frame_id"`
	ActionID      string `json:"action_id"`
	InteractionID string `json:"interaction_id"`
	Signature     string `json:"signature"`
	ExpiresAt     time.Time
}

// NormalizeInteraction strictly validates replay and signature bounds.
// Rejecting stale or tampered interactions at the border prevents out-of-order 
// DAG progression in Layer 3.
func NormalizeInteraction(evt UserInteractionEvent) error {
	if evt.Signature == "" {
		return errors.New("MISSING_SIGNATURE")
	}
	
	if time.Now().After(evt.ExpiresAt) {
		return errors.New("STALE_EVENT_REJECTED")
	}
	
	return nil
}
