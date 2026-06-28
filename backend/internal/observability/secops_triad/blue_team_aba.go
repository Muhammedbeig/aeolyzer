package secopstriad

import (
	"errors"
)

// BehavioralEvent is a sanitized action observation.
type BehavioralEvent struct {
	ActionClass       string  `json:"action_class"`
	ConnectorClass    string  `json:"connector_class,omitempty"`
	PolicyBlocked     bool    `json:"policy_blocked"`
	ApprovalMismatch  bool    `json:"approval_mismatch"`
	CrossTenantSignal bool    `json:"cross_tenant_signal"`
	CostUSD           float64 `json:"cost_usd"`
}

// BehavioralPolicy defines bounded Agent Behavioural Analytics thresholds.
type BehavioralPolicy struct {
	AllowedActionClasses    []string
	AllowedConnectorClasses []string
	MaxPolicyBlocks         int
	MaxCostUSD              float64
}

// BehavioralFinding is a safe, categorical anomaly.
type BehavioralFinding struct {
	Class    string `json:"class"`
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

// AnalyzeBehavior detects expansion, repeated blocks, approval misuse,
// cross-tenant signals, and cost anomalies.
func AnalyzeBehavior(
	events []BehavioralEvent,
	policy BehavioralPolicy,
) ([]BehavioralFinding, error) {
	if len(events) == 0 ||
		len(events) > 10_000 ||
		len(policy.AllowedActionClasses) == 0 ||
		policy.MaxPolicyBlocks < 0 ||
		policy.MaxCostUSD <= 0 {
		return nil, errors.New("behavioral analytics input is invalid")
	}
	allowedActions := set(policy.AllowedActionClasses)
	allowedConnectors := set(policy.AllowedConnectorClasses)
	counts := make(map[string]int)
	totalCost := 0.0
	for _, event := range events {
		if event.ActionClass == "" || event.CostUSD < 0 {
			return nil, errors.New("behavioral event is invalid")
		}
		if _, allowed := allowedActions[event.ActionClass]; !allowed {
			counts["action_class_expansion"]++
		}
		if event.ConnectorClass != "" {
			if _, allowed := allowedConnectors[event.ConnectorClass]; !allowed {
				counts["connector_anomaly"]++
			}
		}
		if event.PolicyBlocked {
			counts["policy_block"]++
		}
		if event.ApprovalMismatch {
			counts["approval_mismatch"]++
		}
		if event.CrossTenantSignal {
			counts["cross_tenant_signal"]++
		}
		totalCost += event.CostUSD
	}
	if counts["policy_block"] <= policy.MaxPolicyBlocks {
		delete(counts, "policy_block")
	}
	if totalCost > policy.MaxCostUSD {
		counts["cost_spike"] = 1
	}
	order := []string{
		"cross_tenant_signal",
		"approval_mismatch",
		"action_class_expansion",
		"connector_anomaly",
		"policy_block",
		"cost_spike",
	}
	var findings []BehavioralFinding
	for _, class := range order {
		count := counts[class]
		if count == 0 {
			continue
		}
		severity := "medium"
		if class == "cross_tenant_signal" || class == "approval_mismatch" {
			severity = "critical"
		} else if class == "action_class_expansion" || class == "connector_anomaly" {
			severity = "high"
		}
		findings = append(findings, BehavioralFinding{
			Class:    class,
			Severity: severity,
			Count:    count,
		})
	}
	return findings, nil
}

func set(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
