package skills

import (
	"testing"
)

func TestLayer4SchemasCompile(t *testing.T) {
	if _, err := NewSchemas(); err != nil {
		t.Fatalf("NewSchemas() failed: %v", err)
	}
}

func TestLayer4SchemasRejectPlaceholders(t *testing.T) {
	schemas, err := NewSchemas()
	if err != nil {
		t.Fatalf("NewSchemas() failed: %v", err)
	}

	tests := map[string]func([]byte) error{
		"registry":          schemas.ValidateRegistry,
		"skill":             schemas.ValidateSkillFrontmatter,
		"resource manifest": schemas.ValidateResourceManifest,
		"eval manifest":     schemas.ValidateEvalManifest,
	}
	for name, validate := range tests {
		t.Run(name, func(t *testing.T) {
			if err := validate([]byte("{}")); err == nil {
				t.Fatal("placeholder document unexpectedly validated")
			}
		})
	}
}

func TestLayer4SchemasAcceptRepresentativeDocuments(t *testing.T) {
	schemas, err := NewSchemas()
	if err != nil {
		t.Fatalf("NewSchemas() failed: %v", err)
	}

	const checksum = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	registry := `version: 2
policy_mode: fail_closed
registry_owner_layer: layer_4_skills
default_status: active
metadata_token_budget:
  max_total_registry_tokens: 2500
  max_description_tokens_per_skill: 90
  max_antitrigger_tokens_per_skill: 60
skills:
  - skill_id: topic_discovery
    name: topic-discovery
    directory: skills/topic_discovery
    status: active
    version: 1.0.0
    owner_team: content_strategy
    tier: read
    risk_class: low
    description: Finds defensible topic candidates. Use for topic planning. Do NOT use for drafting.
    anti_triggers: [drafting requested]
    compatible_profiles: [content_collaborator]
    compatible_intents: [topic_discovery]
    allowed_modes: [plan]
    capability_tags: [topic_discovery]
    declared_action_classes: [read_brand_context]
    output_contracts: [topic_options]
    body_token_estimate: 800
    resource_token_estimate: 0
    eval_manifest: skills/topic_discovery/eval-manifest.yaml
    resource_manifest: skills/topic_discovery/resource-manifest.yaml
    checksum: ` + checksum + "\n"
	if err := schemas.ValidateRegistry([]byte(registry)); err != nil {
		t.Fatalf("ValidateRegistry() failed: %v", err)
	}

	frontmatter := `name: topic-discovery
description: Finds defensible topic candidates. Use for topic planning. Do NOT use for drafting.
version: 1.0.0
owner_team: content_strategy
tier: read
risk_class: low
compatible_profiles: [content_collaborator]
compatible_intents: [topic_discovery]
allowed_modes: [plan]
capability_tags: [topic_discovery]
declared_action_classes: [read_brand_context]
output_contracts: [topic_options]
token_budget:
  body_max_tokens: 1000
  references_max_tokens: 0
  assets_max_tokens: 0
  total_active_max_tokens: 1000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
`
	if err := schemas.ValidateSkillFrontmatter([]byte(frontmatter)); err != nil {
		t.Fatalf("ValidateSkillFrontmatter() failed: %v", err)
	}

	resourceManifest := `version: 1
skill_id: topic_discovery
resources:
  references: []
  assets: []
  scripts: []
`
	if err := schemas.ValidateResourceManifest([]byte(resourceManifest)); err != nil {
		t.Fatalf("ValidateResourceManifest() failed: %v", err)
	}

	evalManifest := `version: 1
skill_id: topic_discovery
eval_owner_layer: layer_8_observability
stored_by_layer: layer_4_skills
minimum_gate:
  trigger_accuracy: 0.9
  negative_trigger_precision: 0.9
  output_rubric_min_score: 4
  trajectory_mode:
    read: ANY_ORDER
    draft: IN_ORDER
    act: EXACT
  regression_pass_required: true
  token_budget_pass_required: true
eval_files:
  trigger_cases: evals/trigger_cases.yaml
  golden_cases: evals/golden_cases.yaml
  trajectory_cases: evals/trajectory_cases.yaml
  rubric: evals/rubric.yaml
  regression_cases: evals/regression_cases.yaml
`
	if err := schemas.ValidateEvalManifest([]byte(evalManifest)); err != nil {
		t.Fatalf("ValidateEvalManifest() failed: %v", err)
	}
}
