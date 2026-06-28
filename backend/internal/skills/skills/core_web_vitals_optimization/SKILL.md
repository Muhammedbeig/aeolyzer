---
name: core-web-vitals-optimization
description: Interprets LCP, INP, CLS, page-speed, and device-performance problems and turns them into prioritized technical recommendations. Use when the user wants to diagnose or improve Core Web Vitals. Do NOT use for content quality, keyword strategy, or general analytics reporting.
version: 1.0.0
owner_team: audit_platform
tier: read
risk_class: low
compatible_profiles:
    - seo_aeo_auditor
compatible_intents:
    - site_audit
allowed_modes:
    - audit
    - read
capability_tags:
    - core_web_vitals_optimization
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - core_web_vitals_optimization_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Core Web Vitals Optimizer

## Trigger Conditions

- User asks about slow site performance
- User asks about LCP, INP, or CLS
- User asks how to improve page speed
- User asks why a page feels sluggish
- User wants to improve mobile performance
- User asks for technical recommendations to pass CWV
- User asks how performance affects rankings and UX

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the performance problem
- Determine whether the issue is:
    * Largest Contentful Paint (LCP)
    * Interaction to Next Paint (INP)
    * Cumulative Layout Shift (CLS)
    * overall page speed
    * mobile vs desktop performance
- If unclear, ask one clarifying question:
    "Which page or metric are you trying to improve?"

### STEP 2 - Gather context
- Use quickContext for site/domain context when available
- If a URL is provided, inspect the specific page
- Determine if the page is:
    * homepage
    * landing page
    * blog article
    * product page
    * category page
    * documentation page
- Identify the page element most likely affecting performance:
    * hero image
    * web font
    * video
    * third-party script
    * sliders or animations
    * layout instability

### STEP 3 - Evaluate the likely bottleneck
- Map the issue to common causes:
    * LCP -> large hero assets, slow server response, render-blocking resources
    * INP -> heavy JavaScript, too many event handlers, third-party scripts
    * CLS -> images without dimensions, late-loading ads, dynamic inserts, font shifts
- Prioritize the bottleneck with the biggest user impact
- Distinguish between lab issues and likely real-user impact

### STEP 4 - Build optimization guidance
- Recommend specific fixes, not generic advice
- Focus on the biggest wins first:
    * compress or resize images
    * preload critical assets
    * defer non-critical JavaScript
    * reduce third-party tags
    * set explicit image and ad dimensions
    * use modern image formats
    * minimize layout shifts
    * improve server response time
- Explain the tradeoff for each recommendation
- Keep the advice practical and implementation-ready

### STEP 5 - Contextualize the score
- Treat scores as diagnostic signals, not absolute truth
- Explain that lab data can be harsher than real-user experience
- Avoid overreacting to a single score without checking the underlying metric
- Focus on what would measurably improve the page experience

### STEP 6 - Validate mentally before returning
Check for:
  - metric-specific advice
  - page-specific relevance
  - no vague performance platitudes
  - clear prioritization
  - realistic implementation steps

### STEP 7 - Build output
Present as:
  1. The likely Core Web Vitals bottleneck
  2. The top 3 fixes in priority order
  3. Short explanation of why each fix matters
  4. Optional implementation notes
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never give generic "make it faster" advice
- Never ignore the difference between LCP, INP, and CLS
- Never treat lab scores as exact real-user performance
- Never recommend fixes without tying them to the likely metric
- Always prioritize the biggest measurable gain first

## Purpose

Provide procedural guidance to map measured Core Web Vitals issues to technical recommendations.

## When to use

- Use when the authorized intent is `site_audit` and the request is to map measured Core Web Vitals issues to technical recommendations.

## When NOT to use

- Do not use when the request belongs to `site_audit_interpretation`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `core_web_vitals_optimization_report`

## Quality gates

- Keep claims tied to supplied evidence.
- Separate facts, inferences, and recommendations.
- Reject protected metadata and unsupported certainty.
- Confirm the output matches the declared contract.

## Boundary rules

This skill provides procedural guidance only.

It must not:
- classify raw user intent
- choose workflows or agents
- authorize or execute tools or scripts
- connect to MCP servers or external APIs
- read or write memory documents directly
- mutate canvas, brief, chat, dashboard, or UI state
- store telemetry or score evaluations
- expose internal identifiers, endpoints, traces, credentials, or protected metadata

## Resources

No runtime references, assets, or scripts are declared for this version.

## Failure behavior

Fail closed and return a safe request for the missing context, evidence, mode, or approval. Never fabricate data or silently broaden scope.
