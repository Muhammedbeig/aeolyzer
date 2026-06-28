package observabilityconfig

import (
	"strings"
	"testing"
)

func TestEmbeddedPoliciesValidate(t *testing.T) {
	policies, err := LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
	}
	if policies.Eval.PassK.HighImpactWorkflow < 10 {
		t.Fatalf(
			"high-impact pass k = %d, want at least 10",
			policies.Eval.PassK.HighImpactWorkflow,
		)
	}
}

func TestDecodeStrictRejectsUnknownFields(t *testing.T) {
	var policy TelemetryPolicy
	err := decodeStrict([]byte("version: 2\npolicy_mode: fail_closed\nunknown: true\n"), &policy)
	if err == nil {
		t.Fatal("decodeStrict() returned nil error")
	}
}

func TestPoliciesRejectUnsafeOverrides(t *testing.T) {
	tests := map[string]func(*Policies){
		"hidden chain of thought": func(p *Policies) {
			p.Telemetry.RedactedReasoning.StoreHiddenChainOfThought = true
		},
		"production red team": func(p *Policies) {
			p.SecOps.RedTeam.RunInProduction = true
		},
		"direct green team recovery": func(p *Policies) {
			p.SecOps.GreenTeam.ExecuteRecoveryDirectly = true
		},
		"weak pass k": func(p *Policies) {
			p.Eval.PassK.HighImpactWorkflow = 1
		},
		"raw protected retention": func(p *Policies) {
			p.Retention.Controls.RawProtectedPayloadRetentionDays = 1
		},
		"mutable governance ledger": func(p *Policies) {
			p.Governance.Ledger.AppendOnly = false
		},
	}

	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			policies, err := LoadEmbeddedPolicies()
			if err != nil {
				t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
			}
			mutate(&policies)
			if err := policies.Validate(); err == nil {
				t.Fatal("Policies.Validate() returned nil error")
			}
		})
	}
}

func TestRedactionPolicyContainsRequiredSecretClasses(t *testing.T) {
	policies, err := LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
	}
	joined := strings.Join(policies.Redaction.NeverStore, ",")
	for _, field := range []string{"raw_api_key", "raw_oauth_token", "private_key", "password"} {
		if !strings.Contains(joined, field) {
			t.Errorf("redaction never_store is missing %q", field)
		}
	}
}
