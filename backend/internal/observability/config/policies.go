// Package observabilityconfig loads and validates Layer 8 policy files.
package observabilityconfig

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"

	"go.yaml.in/yaml/v3"
)

var (
	//go:embed telemetry-policy.yaml
	telemetryPolicyYAML []byte
	//go:embed redaction-policy.yaml
	redactionPolicyYAML []byte
	//go:embed eval-policy.yaml
	evalPolicyYAML []byte
	//go:embed secops-policy.yaml
	secopsPolicyYAML []byte
	//go:embed drift-policy.yaml
	driftPolicyYAML []byte
	//go:embed trust-policy.yaml
	trustPolicyYAML []byte
	//go:embed retention-policy.yaml
	retentionPolicyYAML []byte
	//go:embed governance-policy.yaml
	governancePolicyYAML []byte
)

// Policies contains validated Layer 8 configuration.
type Policies struct {
	Telemetry  TelemetryPolicy
	Redaction  RedactionPolicy
	Eval       EvalPolicy
	SecOps     SecOpsPolicy
	Drift      DriftPolicy
	Trust      TrustPolicy
	Retention  RetentionPolicy
	Governance GovernancePolicy
}

// TelemetryPolicy controls trace completeness, redaction, sampling, and limits.
type TelemetryPolicy struct {
	Version           int    `yaml:"version"`
	PolicyMode        string `yaml:"policy_mode"`
	TraceRequirements struct {
		RequiredRootSpan          string   `yaml:"required_root_span"`
		RequiredChildSpans        []string `yaml:"required_child_spans"`
		AllowMissingOptionalSpans bool     `yaml:"allow_missing_optional_spans"`
		RejectUnlinkedSpans       bool     `yaml:"reject_unlinked_spans"`
		RejectTraceWithoutTenant  bool     `yaml:"reject_trace_without_tenant"`
		RejectTraceWithoutIntent  bool     `yaml:"reject_trace_without_intent"`
		RejectTraceWithoutMode    bool     `yaml:"reject_trace_without_mode"`
	} `yaml:"trace_requirements"`
	RedactedReasoning struct {
		StoreHiddenChainOfThought bool `yaml:"store_hidden_chain_of_thought"`
		StoreReasoningSummary     bool `yaml:"store_reasoning_summary"`
		MaxSummaryChars           int  `yaml:"max_summary_chars"`
		RequireSummaryRedaction   bool `yaml:"require_summary_redaction"`
	} `yaml:"redacted_reasoning"`
	Sampling struct {
		Default                      string  `yaml:"default"`
		RetainAllErrors              bool    `yaml:"retain_all_errors"`
		RetainAllPolicyBlocks        bool    `yaml:"retain_all_policy_blocks"`
		RetainAllQuarantineDecisions bool    `yaml:"retain_all_quarantine_decisions"`
		RetainAllEvalFailures        bool    `yaml:"retain_all_eval_failures"`
		RetainAllHighCostSessions    bool    `yaml:"retain_all_high_cost_sessions"`
		RetainSuccessSampleRate      float64 `yaml:"retain_success_sample_rate"`
	} `yaml:"sampling"`
	CostLatency struct {
		TrackTokenInput             bool    `yaml:"track_token_input"`
		TrackTokenOutput            bool    `yaml:"track_token_output"`
		TrackToolLatencyMS          bool    `yaml:"track_tool_latency_ms"`
		TrackConnectorLatencyMS     bool    `yaml:"track_connector_latency_ms"`
		TrackRuntimeLatencyMS       bool    `yaml:"track_runtime_latency_ms"`
		TrackCostUSD                bool    `yaml:"track_cost_usd"`
		MaxCostPerSessionWarningUSD float64 `yaml:"max_cost_per_session_warning_usd"`
		MaxTurnsWarning             int     `yaml:"max_turns_warning"`
		MaxToolCallsWarning         int     `yaml:"max_tool_calls_warning"`
	} `yaml:"cost_latency"`
}

// RedactionPolicy defines prohibited, hashed, summarized, and protected fields.
type RedactionPolicy struct {
	Version                   int      `yaml:"version"`
	PolicyMode                string   `yaml:"policy_mode"`
	NeverStore                []string `yaml:"never_store"`
	StoreAsHash               []string `yaml:"store_as_hash"`
	StoreAsRedactedSummary    []string `yaml:"store_as_redacted_summary"`
	ProtectedInternalMetadata []string `yaml:"protected_internal_metadata"`
}

// EvalPolicy defines trajectory, pass-k, judge, content, and SEO/AEO gates.
type EvalPolicy struct {
	Version         int    `yaml:"version"`
	PolicyMode      string `yaml:"policy_mode"`
	TrajectoryModes map[string]struct {
		Description string   `yaml:"description"`
		RequiredFor []string `yaml:"required_for"`
		AllowedFor  []string `yaml:"allowed_for"`
	} `yaml:"trajectory_modes"`
	PassK struct {
		Smoke                int `yaml:"smoke"`
		ReadWorkflow         int `yaml:"read_workflow"`
		DraftWorkflow        int `yaml:"draft_workflow"`
		GuardedWriteWorkflow int `yaml:"guarded_write_workflow"`
		HighImpactWorkflow   int `yaml:"high_impact_workflow"`
	} `yaml:"pass_k"`
	Judge struct {
		ScoreMin              int     `yaml:"score_min"`
		ScoreMax              int     `yaml:"score_max"`
		PassThreshold         float64 `yaml:"pass_threshold"`
		UsePositionSwap       bool    `yaml:"use_position_swap"`
		RequireJSONOutput     bool    `yaml:"require_json_output"`
		RequireRubricID       bool    `yaml:"require_rubric_id"`
		RejectWithoutRubric   bool    `yaml:"reject_without_rubric"`
		Temperature           float64 `yaml:"temperature"`
		MinimumConfidence     float64 `yaml:"minimum_confidence"`
		HumanReviewConfidence float64 `yaml:"human_review_confidence"`
		MaxRetries            int     `yaml:"max_retries"`
	} `yaml:"judge"`
	ContentAgentEval struct {
		RequireSourceGroundingForResearch  bool `yaml:"require_source_grounding_for_research"`
		RequireNoSilentMemoryUpdate        bool `yaml:"require_no_silent_memory_update"`
		RequireSelectedTextBindingForEdits bool `yaml:"require_selected_text_binding_for_edits"`
		RequireWordCountTolerance          bool `yaml:"require_word_count_tolerance"`
		WordCountTolerancePercent          int  `yaml:"word_count_tolerance_percent"`
		RequireToneAlignmentScore          bool `yaml:"require_tone_alignment_score"`
	} `yaml:"content_agent_eval"`
	SEOAEOEval struct {
		RequireEvidenceForRecommendations bool `yaml:"require_evidence_for_recommendations"`
		RequireSeverityRationale          bool `yaml:"require_severity_rationale"`
		RequireNoUnsupportedSchemaClaims  bool `yaml:"require_no_unsupported_schema_claims"`
		RequireInternalLinkVerification   bool `yaml:"require_internal_link_verification"`
		RequireCWVMetricMapping           bool `yaml:"require_cwv_metric_mapping"`
	} `yaml:"seo_aeo_eval"`
}

// SecOpsPolicy defines bounded Red, Blue, and Green Team behavior.
type SecOpsPolicy struct {
	Version    int    `yaml:"version"`
	PolicyMode string `yaml:"policy_mode"`
	BlueTeam   struct {
		EnableAgentBehaviourAnalytics bool `yaml:"enable_agent_behaviour_analytics"`
		TrackIntentDrift              bool `yaml:"track_intent_drift"`
		TrackToolAnomalies            bool `yaml:"track_tool_anomalies"`
		TrackConnectorAnomalies       bool `yaml:"track_connector_anomalies"`
		TrackLoopPatterns             bool `yaml:"track_loop_patterns"`
		TrackCostSpikes               bool `yaml:"track_cost_spikes"`
		TrackApprovalMisuse           bool `yaml:"track_approval_misuse"`
		TrackCrossTenantSignals       bool `yaml:"track_cross_tenant_signals"`
	} `yaml:"blue_team"`
	RedTeam struct {
		RunInProduction bool     `yaml:"run_in_production"`
		RunInShadow     bool     `yaml:"run_in_shadow"`
		RunInCI         bool     `yaml:"run_in_ci"`
		PayloadClasses  []string `yaml:"payload_classes"`
	} `yaml:"red_team"`
	GreenTeam struct {
		ExecuteRecoveryDirectly             bool   `yaml:"execute_recovery_directly"`
		PublishRecoveryRecommendations      bool   `yaml:"publish_recovery_recommendations"`
		PublishQuarantineDecisions          bool   `yaml:"publish_quarantine_decisions"`
		QuarantineExecutorLayer             string `yaml:"quarantine_executor_layer"`
		RequireHumanReviewForPolicyChange   bool   `yaml:"require_human_review_for_policy_change"`
		RequireHumanReviewForSkillChange    bool   `yaml:"require_human_review_for_skill_change"`
		RequireHumanReviewForWorkflowChange bool   `yaml:"require_human_review_for_workflow_change"`
	} `yaml:"green_team"`
	SmallBatchEnforcement struct {
		MaxFilesChangedBeforeHighRisk   int `yaml:"max_files_changed_before_high_risk"`
		MaxSurfaceMutationsBeforeReview int `yaml:"max_surface_mutations_before_review"`
		MaxConsecutiveToolFailures      int `yaml:"max_consecutive_tool_failures"`
		MaxUnresolvedPolicyBlocks       int `yaml:"max_unresolved_policy_blocks"`
	} `yaml:"small_batch_enforcement"`
}

// DriftPolicy defines drift, trust-decay, recovery, and loop thresholds.
type DriftPolicy struct {
	Version     int    `yaml:"version"`
	PolicyMode  string `yaml:"policy_mode"`
	IntentDrift struct {
		CompareAgainst []string `yaml:"compare_against"`
		Thresholds     struct {
			Warn                float64 `yaml:"warn"`
			BlockRecommendation float64 `yaml:"block_recommendation"`
			QuarantineDecision  float64 `yaml:"quarantine_decision"`
		} `yaml:"thresholds"`
	} `yaml:"intent_drift"`
	TrustDecay struct {
		StartingScore    float64            `yaml:"starting_score"`
		MinimumSafeScore float64            `yaml:"minimum_safe_score"`
		QuarantineScore  float64            `yaml:"quarantine_score"`
		DecayFactors     map[string]float64 `yaml:"decay_factors"`
		RecoveryFactors  map[string]float64 `yaml:"recovery_factors"`
	} `yaml:"trust_decay"`
	Loops struct {
		MaxRepeatedNode           int `yaml:"max_repeated_node"`
		MaxRepeatedToolSameParams int `yaml:"max_repeated_tool_same_params"`
		MaxReplanCount            int `yaml:"max_replan_count"`
	} `yaml:"loops"`
}

// TrustPolicy defines score bounds and critical-signal behavior.
type TrustPolicy struct {
	Version    int    `yaml:"version"`
	PolicyMode string `yaml:"policy_mode"`
	Score      struct {
		Minimum             float64 `yaml:"minimum"`
		Maximum             float64 `yaml:"maximum"`
		Initial             float64 `yaml:"initial"`
		Warning             float64 `yaml:"warning"`
		QuarantineDecision  float64 `yaml:"quarantine_decision"`
		RecoveryCapPerTrace float64 `yaml:"recovery_cap_per_trace"`
	} `yaml:"score"`
	CriticalSignals []string `yaml:"critical_signals"`
	Rules           struct {
		CriticalSignalRequiresQuarantineDecision bool `yaml:"critical_signal_requires_quarantine_decision"`
		ScoreBelowQuarantineRequiresDecision     bool `yaml:"score_below_quarantine_requires_decision"`
		Layer8MustNotExecuteQuarantine           bool `yaml:"layer8_must_not_execute_quarantine"`
		RequirePolicyVersionInDecision           bool `yaml:"require_policy_version_in_decision"`
		RequireEvidenceRefsInDecision            bool `yaml:"require_evidence_refs_in_decision"`
	} `yaml:"rules"`
}

// RetentionPolicy defines retention classes and privacy controls.
type RetentionPolicy struct {
	Version    int    `yaml:"version"`
	PolicyMode string `yaml:"policy_mode"`
	Classes    map[string]struct {
		RetentionDays int `yaml:"retention_days"`
	} `yaml:"classes"`
	Controls struct {
		LegalHoldSupported               bool `yaml:"legal_hold_supported"`
		TenantDeletionSupported          bool `yaml:"tenant_deletion_supported"`
		EncryptionAtRestRequired         bool `yaml:"encryption_at_rest_required"`
		ImmutableGovernanceRecords       bool `yaml:"immutable_governance_records"`
		RawProtectedPayloadRetentionDays int  `yaml:"raw_protected_payload_retention_days"`
		DeletionAuditRequired            bool `yaml:"deletion_audit_required"`
	} `yaml:"controls"`
}

// GovernancePolicy defines immutable-ledger and release requirements.
type GovernancePolicy struct {
	Version    int    `yaml:"version"`
	PolicyMode string `yaml:"policy_mode"`
	Ledger     struct {
		AppendOnly          bool `yaml:"append_only"`
		HashChainRequired   bool `yaml:"hash_chain_required"`
		SignerRequired      bool `yaml:"signer_required"`
		TimestampRequired   bool `yaml:"timestamp_required"`
		TenantScopeRequired bool `yaml:"tenant_scope_required"`
	} `yaml:"ledger"`
	HumanReviewRequiredFor []string `yaml:"human_review_required_for"`
	Release                struct {
		RequireAllSafetyEvals         bool `yaml:"require_all_safety_evals"`
		RequirePassK                  bool `yaml:"require_pass_k"`
		RequireNoOpenCriticalFindings bool `yaml:"require_no_open_critical_findings"`
		RequireCanaryEvidence         bool `yaml:"require_canary_evidence"`
		RequireRollbackPlan           bool `yaml:"require_rollback_plan"`
		RequireNamedApprover          bool `yaml:"require_named_approver"`
	} `yaml:"release"`
}

// LoadEmbeddedPolicies strictly parses and semantically validates every Layer
// 8 policy. It fails closed on unknown fields.
func LoadEmbeddedPolicies() (Policies, error) {
	var policies Policies
	loaders := []struct {
		name string
		data []byte
		dst  any
	}{
		{name: "telemetry", data: telemetryPolicyYAML, dst: &policies.Telemetry},
		{name: "redaction", data: redactionPolicyYAML, dst: &policies.Redaction},
		{name: "eval", data: evalPolicyYAML, dst: &policies.Eval},
		{name: "secops", data: secopsPolicyYAML, dst: &policies.SecOps},
		{name: "drift", data: driftPolicyYAML, dst: &policies.Drift},
		{name: "trust", data: trustPolicyYAML, dst: &policies.Trust},
		{name: "retention", data: retentionPolicyYAML, dst: &policies.Retention},
		{name: "governance", data: governancePolicyYAML, dst: &policies.Governance},
	}
	for _, loader := range loaders {
		if err := decodeStrict(loader.data, loader.dst); err != nil {
			return Policies{}, fmt.Errorf("decode %s policy: %w", loader.name, err)
		}
	}
	if err := policies.Validate(); err != nil {
		return Policies{}, err
	}
	return policies, nil
}

// Validate checks cross-field security invariants in all Layer 8 policies.
func (p Policies) Validate() error {
	for name, header := range map[string]struct {
		version int
		mode    string
	}{
		"telemetry":  {p.Telemetry.Version, p.Telemetry.PolicyMode},
		"redaction":  {p.Redaction.Version, p.Redaction.PolicyMode},
		"eval":       {p.Eval.Version, p.Eval.PolicyMode},
		"secops":     {p.SecOps.Version, p.SecOps.PolicyMode},
		"drift":      {p.Drift.Version, p.Drift.PolicyMode},
		"trust":      {p.Trust.Version, p.Trust.PolicyMode},
		"retention":  {p.Retention.Version, p.Retention.PolicyMode},
		"governance": {p.Governance.Version, p.Governance.PolicyMode},
	} {
		if header.version != 2 || header.mode != "fail_closed" {
			return fmt.Errorf("%s policy must be version 2 and fail closed", name)
		}
	}
	if err := p.validateTelemetry(); err != nil {
		return err
	}
	if err := p.validateRedaction(); err != nil {
		return err
	}
	if err := p.validateEval(); err != nil {
		return err
	}
	if err := p.validateSecOps(); err != nil {
		return err
	}
	if err := p.validateDriftAndTrust(); err != nil {
		return err
	}
	if err := p.validateRetentionAndGovernance(); err != nil {
		return err
	}
	return nil
}

func (p Policies) validateTelemetry() error {
	t := p.Telemetry
	if t.TraceRequirements.RequiredRootSpan == "" ||
		len(t.TraceRequirements.RequiredChildSpans) == 0 ||
		!t.TraceRequirements.RejectUnlinkedSpans ||
		!t.TraceRequirements.RejectTraceWithoutTenant ||
		!t.TraceRequirements.RejectTraceWithoutIntent ||
		!t.TraceRequirements.RejectTraceWithoutMode {
		return errors.New("telemetry trace requirements are not fail closed")
	}
	if t.RedactedReasoning.StoreHiddenChainOfThought ||
		!t.RedactedReasoning.StoreReasoningSummary ||
		!t.RedactedReasoning.RequireSummaryRedaction ||
		t.RedactedReasoning.MaxSummaryChars < 1 {
		return errors.New("telemetry reasoning policy is unsafe")
	}
	if t.Sampling.Default != "tail_based" ||
		t.Sampling.RetainSuccessSampleRate < 0 ||
		t.Sampling.RetainSuccessSampleRate > 1 ||
		!t.Sampling.RetainAllErrors ||
		!t.Sampling.RetainAllPolicyBlocks ||
		!t.Sampling.RetainAllQuarantineDecisions ||
		!t.Sampling.RetainAllEvalFailures {
		return errors.New("telemetry sampling policy is unsafe")
	}
	if t.CostLatency.MaxCostPerSessionWarningUSD <= 0 ||
		t.CostLatency.MaxTurnsWarning < 1 ||
		t.CostLatency.MaxToolCallsWarning < 1 {
		return errors.New("telemetry cost and latency limits are invalid")
	}
	return nil
}

func (p Policies) validateRedaction() error {
	required := []string{
		"hidden_chain_of_thought",
		"raw_system_prompt",
		"raw_developer_prompt",
		"raw_user_prompt",
		"raw_api_key",
		"raw_oauth_token",
		"password",
		"private_key",
		"raw_pii",
	}
	neverStore := stringSet(p.Redaction.NeverStore)
	for _, field := range required {
		if _, found := neverStore[field]; !found {
			return fmt.Errorf("redaction policy must prohibit %s", field)
		}
	}
	if overlaps(p.Redaction.NeverStore, p.Redaction.StoreAsHash) ||
		overlaps(p.Redaction.NeverStore, p.Redaction.StoreAsRedactedSummary) {
		return errors.New("redaction storage classes overlap")
	}
	return nil
}

func (p Policies) validateEval() error {
	for _, mode := range []string{"EXACT", "IN_ORDER", "ANY_ORDER"} {
		if _, found := p.Eval.TrajectoryModes[mode]; !found {
			return fmt.Errorf("eval policy is missing %s trajectory mode", mode)
		}
	}
	k := p.Eval.PassK
	if k.Smoke != 1 ||
		k.ReadWorkflow < 3 ||
		k.DraftWorkflow < 3 ||
		k.GuardedWriteWorkflow < 5 ||
		k.HighImpactWorkflow < 10 {
		return errors.New("eval pass k policy is below production minimum")
	}
	j := p.Eval.Judge
	if j.ScoreMin != 1 ||
		j.ScoreMax != 5 ||
		j.PassThreshold < 4 ||
		!j.UsePositionSwap ||
		!j.RequireJSONOutput ||
		!j.RequireRubricID ||
		!j.RejectWithoutRubric ||
		j.Temperature != 0 ||
		j.MinimumConfidence < 0.8 ||
		j.HumanReviewConfidence < j.MinimumConfidence ||
		j.MaxRetries < 0 ||
		j.MaxRetries > 2 {
		return errors.New("eval judge policy is unsafe")
	}
	return nil
}

func (p Policies) validateSecOps() error {
	s := p.SecOps
	if s.RedTeam.RunInProduction ||
		!s.RedTeam.RunInCI ||
		len(s.RedTeam.PayloadClasses) == 0 ||
		s.GreenTeam.ExecuteRecoveryDirectly ||
		s.GreenTeam.QuarantineExecutorLayer != "layer_6_runtime" ||
		!s.GreenTeam.RequireHumanReviewForPolicyChange ||
		!s.GreenTeam.RequireHumanReviewForSkillChange ||
		!s.GreenTeam.RequireHumanReviewForWorkflowChange {
		return errors.New("secops team boundaries are unsafe")
	}
	if s.SmallBatchEnforcement.MaxFilesChangedBeforeHighRisk < 1 ||
		s.SmallBatchEnforcement.MaxSurfaceMutationsBeforeReview < 1 ||
		s.SmallBatchEnforcement.MaxConsecutiveToolFailures < 1 ||
		s.SmallBatchEnforcement.MaxUnresolvedPolicyBlocks < 1 {
		return errors.New("secops small-batch limits are invalid")
	}
	return nil
}

func (p Policies) validateDriftAndTrust() error {
	thresholds := p.Drift.IntentDrift.Thresholds
	if thresholds.Warn <= 0 ||
		thresholds.Warn >= thresholds.BlockRecommendation ||
		thresholds.BlockRecommendation >= thresholds.QuarantineDecision ||
		thresholds.QuarantineDecision > 1 {
		return errors.New("drift thresholds must increase toward quarantine")
	}
	decay := p.Drift.TrustDecay
	if decay.StartingScore != 1 ||
		decay.QuarantineScore < 0 ||
		decay.QuarantineScore >= decay.MinimumSafeScore ||
		decay.MinimumSafeScore > decay.StartingScore ||
		len(decay.DecayFactors) == 0 ||
		len(decay.RecoveryFactors) == 0 {
		return errors.New("drift trust-decay policy is invalid")
	}
	score := p.Trust.Score
	if score.Minimum != 0 ||
		score.Maximum != 1 ||
		score.Initial != 1 ||
		score.QuarantineDecision >= score.Warning ||
		score.RecoveryCapPerTrace <= 0 ||
		score.RecoveryCapPerTrace > 0.25 ||
		len(p.Trust.CriticalSignals) == 0 ||
		!p.Trust.Rules.Layer8MustNotExecuteQuarantine {
		return errors.New("trust policy is invalid")
	}
	return nil
}

func (p Policies) validateRetentionAndGovernance() error {
	if len(p.Retention.Classes) == 0 ||
		p.Retention.Controls.RawProtectedPayloadRetentionDays != 0 ||
		!p.Retention.Controls.TenantDeletionSupported ||
		!p.Retention.Controls.EncryptionAtRestRequired ||
		!p.Retention.Controls.ImmutableGovernanceRecords ||
		!p.Retention.Controls.DeletionAuditRequired {
		return errors.New("retention policy is unsafe")
	}
	g := p.Governance
	if !g.Ledger.AppendOnly ||
		!g.Ledger.HashChainRequired ||
		!g.Ledger.SignerRequired ||
		!g.Ledger.TimestampRequired ||
		!g.Ledger.TenantScopeRequired ||
		len(g.HumanReviewRequiredFor) == 0 ||
		!g.Release.RequireAllSafetyEvals ||
		!g.Release.RequirePassK ||
		!g.Release.RequireNoOpenCriticalFindings ||
		!g.Release.RequireCanaryEvidence ||
		!g.Release.RequireRollbackPlan ||
		!g.Release.RequireNamedApprover {
		return errors.New("governance policy is unsafe")
	}
	return nil
}

func decodeStrict(data []byte, dst any) error {
	if len(data) == 0 {
		return errors.New("policy document is empty")
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("policy contains multiple yaml documents")
	}
	return nil
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func overlaps(left, right []string) bool {
	seen := stringSet(left)
	for _, value := range right {
		if _, found := seen[value]; found {
			return true
		}
	}
	return false
}
