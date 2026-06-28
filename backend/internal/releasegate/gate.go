// Package releasegate verifies repository evidence required for a production release.
//
// The gate does not implement any agent layer. It checks whether the artifacts and
// executable controls required by the repository specifications are present and
// non-placeholder before a release can be described as production-ready.
package releasegate

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Finding describes one production-readiness blocker.
type Finding struct {
	Code    string `json:"code"`
	Area    string `json:"area"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
}

// Report contains all blockers found by Check.
type Report struct {
	Findings []Finding `json:"findings"`
}

// Ready reports whether the repository passed every release gate.
func (r Report) Ready() bool {
	return len(r.Findings) == 0
}

// Check evaluates production evidence rooted at root.
func Check(root string) (Report, error) {
	info, err := os.Stat(root)
	if err != nil {
		return Report{}, fmt.Errorf("stat repository root: %w", err)
	}
	if !info.IsDir() {
		return Report{}, errors.New("repository root is not a directory")
	}

	checker := &repositoryChecker{root: root}
	checker.checkPlaceholderArtifacts()
	checker.checkCriticalControls()
	checker.checkSkills()
	checker.checkExplicitPrototypeMarkers()
	sort.Slice(checker.findings, func(i, j int) bool {
		left := checker.findings[i]
		right := checker.findings[j]
		if left.Area != right.Area {
			return left.Area < right.Area
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		return left.Path < right.Path
	})

	return Report{Findings: checker.findings}, nil
}

type repositoryChecker struct {
	root     string
	findings []Finding
}

type artifactRequirement struct {
	area string
	path string
}

var configuredArtifacts = []artifactRequirement{
	{area: "layer_2", path: "config/policies.yaml"},
	{area: "layer_2", path: "config/routing-schema.json"},
	{area: "layer_4", path: "internal/skills/eval-manifest.schema.json"},
	{area: "layer_4", path: "internal/skills/registry.schema.json"},
	{area: "layer_4", path: "internal/skills/resource-manifest.schema.json"},
	{area: "layer_4", path: "internal/skills/skill.schema.json"},
	{area: "layer_4", path: "internal/skills/skill-registry.yaml"},
	{area: "layer_5", path: "internal/extensions/a2a-agent-card.schema.json"},
	{area: "layer_5", path: "internal/extensions/a2a-envelope.schema.json"},
	{area: "layer_5", path: "internal/extensions/a2ui-catalog.schema.json"},
	{area: "layer_5", path: "internal/extensions/a2ui-frame.schema.json"},
	{area: "layer_5", path: "internal/extensions/approval.schema.json"},
	{area: "layer_5", path: "internal/extensions/catalog-lock.yaml"},
	{area: "layer_5", path: "internal/extensions/presentation.schema.json"},
	{area: "layer_5", path: "internal/extensions/surface-patch.schema.json"},
	{area: "layer_5", path: "internal/extensions/ui-event.schema.json"},
	{area: "layer_6", path: "internal/runtime/dependency-policy.schema.json"},
	{area: "layer_6", path: "internal/runtime/egress-policy.schema.json"},
	{area: "layer_6", path: "internal/runtime/filesystem-policy.schema.json"},
	{area: "layer_6", path: "internal/runtime/jit-token.schema.json"},
	{area: "layer_6", path: "internal/runtime/quarantine-command.schema.json"},
	{area: "layer_6", path: "internal/runtime/runtime-execution.schema.json"},
	{area: "layer_6", path: "internal/runtime/runtime-result.schema.json"},
	{area: "layer_6", path: "internal/runtime/sandbox-lease.schema.json"},
	{area: "layer_7", path: "internal/interop/config/connector-registry.yaml"},
	{area: "layer_7", path: "internal/interop/config/mcp-server-manifest.schema.json"},
	{area: "layer_7", path: "internal/interop/config/source-contracts.yaml"},
	{area: "layer_8", path: "internal/observability/config/drift-policy.yaml"},
	{area: "layer_8", path: "internal/observability/config/eval-policy.yaml"},
	{area: "layer_8", path: "internal/observability/config/redaction-policy.yaml"},
	{area: "layer_8", path: "internal/observability/config/secops-policy.yaml"},
	{area: "layer_8", path: "internal/observability/config/telemetry-policy.yaml"},
}

var criticalControls = []artifactRequirement{
	{area: "layer_5", path: "internal/extensions/a2ui_translator/a2ui_schema_manager.go"},
	{area: "layer_5", path: "internal/extensions/a2ui_translator/a2ui_part_converter.go"},
	{area: "layer_5", path: "internal/extensions/security/hidden_payload_scanner.go"},
	{area: "layer_5", path: "internal/extensions/security/markdown_sanitizer.go"},
	{area: "layer_5", path: "internal/extensions/security/ui_payload_signer.go"},
	{area: "layer_5", path: "internal/extensions/a2a_server/a2a_envelope_validator.go"},
	{area: "layer_6", path: "internal/runtime/sandbox_environment/gvisor_isolation.go"},
	{area: "layer_6", path: "internal/runtime/sandbox_environment/ephemeral_state_manager.go"},
	{area: "layer_6", path: "internal/runtime/network_egress/egress_proxy_controller.go"},
	{area: "layer_6", path: "internal/runtime/supply_chain_defense/dynamic_install_blocker.go"},
	{area: "layer_6", path: "internal/runtime/iam_context/jit_token_broker.go"},
	{area: "layer_6", path: "internal/runtime/secops_green_team_ops/stateful_quarantine.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/jsonrpc_codec.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/stdio_client.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/streamable_http_client.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/handshake_validator.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/schema_hash_verifier.go"},
	{area: "layer_7", path: "internal/interop/mcp_transport_plane/circuit_breaker.go"},
	{area: "layer_7", path: "internal/interop/data_security_mesh/mtls_identity_verifier.go"},
	{area: "layer_7", path: "internal/interop/data_security_mesh/token_audience_validator.go"},
	{area: "layer_7", path: "internal/interop/data_security_mesh/field_projection_enforcer.go"},
	{area: "layer_7", path: "internal/interop/data_security_mesh/provenance_attestor.go"},
	{area: "layer_7", path: "internal/interop/data_security_mesh/taint_marker.go"},
	{area: "layer_7", path: "internal/interop/vector_rag_store/vector_retrieval_engine.go"},
	{area: "layer_7", path: "internal/interop/vector_rag_store/tenant_partitioning.go"},
	{area: "layer_8", path: "internal/observability/telemetry_tracing/opentelemetry_tracker.go"},
	{area: "layer_8", path: "internal/observability/drift_trust/intent_drift_detector.go"},
	{area: "layer_8", path: "internal/observability/drift_trust/trust_decay_scorer.go"},
	{area: "layer_8", path: "internal/observability/drift_trust/loop_detector.go"},
	{area: "layer_8", path: "internal/observability/secops_triad/red_team_simulator.go"},
	{area: "layer_8", path: "internal/observability/secops_triad/blue_team_aba.go"},
	{area: "layer_8", path: "internal/observability/secops_triad/green_team_recovery_planner.go"},
	{area: "layer_8", path: "internal/observability/evaluation_engine/trajectory_evaluator.go"},
	{area: "layer_8", path: "internal/observability/evaluation_engine/llm_as_judge.go"},
	{area: "layer_8", path: "internal/observability/evaluation_engine/pass_k_runner.go"},
	{area: "layer_8", path: "internal/observability/governance_audit/immutable_audit_ledger.go"},
	{area: "layer_8", path: "internal/observability/feedback_improvement_loop/correction_miner.go"},
}

func (c *repositoryChecker) checkPlaceholderArtifacts() {
	for _, requirement := range configuredArtifacts {
		path := c.fullPath(requirement.path)
		data, err := os.ReadFile(path)
		if err != nil {
			code := "required_artifact_unreadable"
			if errors.Is(err, os.ErrNotExist) {
				code = "required_artifact_missing"
			}
			c.add(code, requirement.area, requirement.path, "required schema or policy is unavailable")
			continue
		}
		if placeholderContent(data) {
			c.add("placeholder_artifact", requirement.area, requirement.path, "schema or policy contains placeholder-only content")
		}
	}
}

func (c *repositoryChecker) checkCriticalControls() {
	for _, requirement := range criticalControls {
		info, err := os.Stat(c.fullPath(requirement.path))
		if err != nil || info.IsDir() || info.Size() == 0 {
			c.add("executable_control_missing", requirement.area, requirement.path, "required executable production control is not implemented")
		}
	}
}

func (c *repositoryChecker) checkSkills() {
	skillsRoot := c.fullPath("internal/skills/skills")
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		c.add("skill_library_unreadable", "layer_4", "internal/skills/skills", "skill library cannot be read")
		return
	}

	var skillNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			skillNames = append(skillNames, entry.Name())
		}
	}
	sort.Strings(skillNames)
	if len(skillNames) == 0 {
		c.add("skill_library_empty", "layer_4", "internal/skills/skills", "no skill directories were found")
		return
	}

	registryPath := c.fullPath("internal/skills/skill-registry.yaml")
	registryData, err := os.ReadFile(registryPath)
	if err != nil {
		c.add("skill_registry_unreadable", "layer_4", "internal/skills/skill-registry.yaml", "skill registry cannot be read")
	} else {
		registryIDs := registrySkillIDs(string(registryData))
		if len(registryIDs) == 0 {
			c.add("skill_registry_empty", "layer_4", "internal/skills/skill-registry.yaml", fmt.Sprintf("%d skill directories exist but none are registered", len(skillNames)))
		} else {
			var unregistered []string
			for _, name := range skillNames {
				if _, ok := registryIDs[name]; !ok {
					unregistered = append(unregistered, name)
				}
			}
			c.addSkillAggregate("skill_unregistered", "skill directories are absent from the registry", unregistered, len(skillNames))
		}
	}

	missingFrontmatter := make(map[string][]string)
	missingSections := make(map[string][]string)
	missingEvidence := make(map[string][]string)
	oversizedBodies := make([]string, 0)

	for _, name := range skillNames {
		relativeRoot := filepath.Join("internal", "skills", "skills", name)
		skillPath := filepath.Join(relativeRoot, "SKILL.md")
		data, err := os.ReadFile(c.fullPath(skillPath))
		if err != nil {
			missingEvidence["SKILL.md"] = append(missingEvidence["SKILL.md"], name)
			continue
		}

		fields, body, ok := parseSkillFile(string(data))
		if !ok {
			missingFrontmatter["valid YAML frontmatter envelope"] = append(missingFrontmatter["valid YAML frontmatter envelope"], name)
		}
		for _, field := range requiredSkillFields {
			if _, exists := fields[field]; !exists {
				missingFrontmatter[field] = append(missingFrontmatter[field], name)
			}
		}
		for _, section := range requiredSkillSections {
			if !hasMarkdownSection(body, section) {
				missingSections[section] = append(missingSections[section], name)
			}
		}
		if len(strings.Fields(body)) > 5000 {
			oversizedBodies = append(oversizedBodies, name)
		}

		for _, evidence := range requiredSkillEvidence {
			path := filepath.Join(relativeRoot, evidence)
			info, err := os.Stat(c.fullPath(path))
			if err != nil || info.IsDir() || info.Size() == 0 {
				missingEvidence[evidence] = append(missingEvidence[evidence], name)
			}
		}
	}

	c.addGroupedSkillFindings("skill_frontmatter_incomplete", "required frontmatter field is missing", missingFrontmatter, len(skillNames))
	c.addGroupedSkillFindings("skill_sections_incomplete", "required SKILL.md section is missing", missingSections, len(skillNames))
	c.addGroupedSkillFindings("skill_eval_evidence_missing", "required ownership, release, resource, or evaluation evidence is missing", missingEvidence, len(skillNames))
	c.addSkillAggregate("skill_body_over_limit", "SKILL.md body exceeds the 5,000-word hard limit", oversizedBodies, len(skillNames))
}

var requiredSkillFields = []string{
	"name",
	"description",
	"version",
	"owner_team",
	"tier",
	"risk_class",
	"compatible_profiles",
	"compatible_intents",
	"allowed_modes",
	"capability_tags",
	"declared_action_classes",
	"output_contracts",
	"token_budget",
	"resource_manifest",
	"eval_manifest",
}

var requiredSkillSections = []string{
	"Purpose",
	"When to use",
	"When NOT to use",
	"Inputs expected",
	"Procedure",
	"Output contract",
	"Quality gates",
	"Boundary rules",
	"Resources",
	"Failure behavior",
}

var requiredSkillEvidence = []string{
	"OWNERS",
	"CHANGELOG.md",
	"resource-manifest.yaml",
	"eval-manifest.yaml",
	"evals/trigger_cases.yaml",
	"evals/golden_cases.yaml",
	"evals/trajectory_cases.yaml",
	"evals/rubric.yaml",
	"evals/regression_cases.yaml",
}

func (c *repositoryChecker) checkExplicitPrototypeMarkers() {
	markers := map[string]string{
		"placeholder":                               "source explicitly identifies placeholder behavior",
		"production path should swap this":          "source explicitly defers its production implementation",
		"vulnerable to adversarial token splitting": "source explicitly documents a known security bypass",
		"should ideally be supplemented":            "source explicitly documents an incomplete security control",
		"dummy test":                                "test explicitly identifies itself as non-verifying",
	}

	roots := []string{"cmd", "internal"}
	for _, root := range roots {
		_ = filepath.WalkDir(c.fullPath(root), func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				return nil
			}
			relative, err := filepath.Rel(c.root, path)
			if err != nil {
				return nil
			}
			relative = filepath.ToSlash(relative)
			if strings.HasPrefix(relative, "internal/releasegate/") ||
				strings.HasPrefix(relative, "cmd/readiness/") {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 0
			for scanner.Scan() {
				lineNumber++
				lower := strings.ToLower(scanner.Text())
				for marker, message := range markers {
					if strings.Contains(lower, marker) {
						c.add("explicit_prototype_marker", "implementation", relative+fmt.Sprintf(":%d", lineNumber), message)
						break
					}
				}
			}
			return nil
		})
	}
}

func (c *repositoryChecker) addGroupedSkillFindings(code, message string, grouped map[string][]string, total int) {
	keys := make([]string, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		names := grouped[key]
		if len(names) == 0 {
			continue
		}
		c.add(code, "layer_4", "", fmt.Sprintf("%s: %s (%d/%d skills; examples: %s)", message, key, len(names), total, examples(names)))
	}
}

func (c *repositoryChecker) addSkillAggregate(code, message string, names []string, total int) {
	if len(names) == 0 {
		return
	}
	c.add(code, "layer_4", "", fmt.Sprintf("%s (%d/%d skills; examples: %s)", message, len(names), total, examples(names)))
}

func (c *repositoryChecker) add(code, area, path, message string) {
	c.findings = append(c.findings, Finding{
		Code:    code,
		Area:    area,
		Path:    filepath.ToSlash(path),
		Message: message,
	})
}

func (c *repositoryChecker) fullPath(path string) string {
	return filepath.Join(c.root, filepath.FromSlash(path))
}

func placeholderContent(data []byte) bool {
	value := strings.TrimSpace(string(data))
	switch value {
	case "", "{}", "[]", "version: 1", "version: 2", "skills: []":
		return true
	default:
		return false
	}
}

func registrySkillIDs(data string) map[string]struct{} {
	result := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "skill_id:") && !strings.HasPrefix(line, "- skill_id:") {
			continue
		}
		_, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		id := strings.Trim(strings.TrimSpace(value), `"'`)
		if id != "" {
			result[strings.ReplaceAll(id, "-", "_")] = struct{}{}
		}
	}
	return result
}

func parseSkillFile(data string) (map[string]struct{}, string, bool) {
	if !strings.HasPrefix(data, "---\n") && !strings.HasPrefix(data, "---\r\n") {
		return nil, data, false
	}
	normalized := strings.ReplaceAll(data, "\r\n", "\n")
	parts := strings.SplitN(normalized, "\n---\n", 2)
	if len(parts) != 2 {
		return nil, data, false
	}

	fields := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimPrefix(parts[0], "---\n")))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == ' ' || line[0] == '\t' {
			continue
		}
		key, _, ok := strings.Cut(line, ":")
		if ok && key != "" {
			fields[strings.TrimSpace(key)] = struct{}{}
		}
	}
	return fields, parts[1], true
}

func hasMarkdownSection(body, section string) bool {
	target := "## " + section
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		if strings.EqualFold(strings.TrimSpace(scanner.Text()), target) {
			return true
		}
	}
	return false
}

func examples(names []string) string {
	const limit = 4
	if len(names) <= limit {
		return strings.Join(names, ", ")
	}
	return strings.Join(names[:limit], ", ") + ", ..."
}
