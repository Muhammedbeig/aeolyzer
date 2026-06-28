package intake_events

import (
	"time"

	"aeolyzer/layer_02_intake/contracts"
)

type SafetyEmitter struct {
	emitFunc func(event contracts.SafetyEvent)
}

func NewSafetyEmitter(emit func(event contracts.SafetyEvent)) *SafetyEmitter {
	return &SafetyEmitter{emitFunc: emit}
}

func (s *SafetyEmitter) Emit(event contracts.SafetyEvent) {
	if s.emitFunc != nil {
		if event.CreatedAt.IsZero() {
			event.CreatedAt = time.Now()
		}
		s.emitFunc(event)
	}
}

func (s *SafetyEmitter) EmitToolPolicyEvent(traceID, toolName, actionClass, decision, reasonCode string) {
	s.Emit(contracts.SafetyEvent{
		TraceID:     traceID,
		EventType:   "tool_policy_" + decision,
		ActionClass: actionClass,
		Decision:    decision,
		ReasonCode:  reasonCode,
		Metadata: map[string]interface{}{
			"tool_name": toolName,
		},
	})
}
