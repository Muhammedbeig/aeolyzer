package secopstriad

import (
	"errors"
	"sort"
)

// RecoveryPlan is a recommendation only. Layer 6 executes allowed quarantine
// actions after signed decision validation.
type RecoveryPlan struct {
	Severity            string   `json:"severity"`
	RecommendedActions  []string `json:"recommended_actions"`
	RequiresHumanReview bool     `json:"requires_human_review"`
	SafeSummary         string   `json:"safe_summary"`
}

// PlanRecovery maps findings to bounded recovery recommendations.
func PlanRecovery(findings []BehavioralFinding) (RecoveryPlan, error) {
	if len(findings) == 0 {
		return RecoveryPlan{}, errors.New("recovery findings are required")
	}
	actionSet := make(map[string]struct{})
	severity := "medium"
	for _, finding := range findings {
		if finding.Class == "" || finding.Count < 1 {
			return RecoveryPlan{}, errors.New("recovery finding is invalid")
		}
		switch finding.Class {
		case "cross_tenant_signal", "approval_mismatch":
			severity = "critical"
			actionSet["stop_new_executions"] = struct{}{}
			actionSet["revoke_jit_tokens"] = struct{}{}
			actionSet["block_egress"] = struct{}{}
			actionSet["preserve_forensic_snapshot"] = struct{}{}
		case "action_class_expansion", "connector_anomaly":
			if severity != "critical" {
				severity = "high"
			}
			actionSet["revoke_tool_access"] = struct{}{}
			actionSet["block_egress"] = struct{}{}
		case "policy_block", "cost_spike":
			actionSet["stop_new_executions"] = struct{}{}
			actionSet["allow_read_only_status"] = struct{}{}
		default:
			return RecoveryPlan{}, errors.New("unsupported recovery finding")
		}
	}
	actions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	return RecoveryPlan{
		Severity:            severity,
		RecommendedActions:  actions,
		RequiresHumanReview: true,
		SafeSummary:         "Agent activity requires bounded containment and human review.",
	}, nil
}
