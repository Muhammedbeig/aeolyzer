---
name: google-business-profile-optimization
description: Optimizes Google Business Profile relevance, prominence, completeness, categories, services, reviews, photos, posts, and Q&A. Use when the user wants stronger Maps or local-pack performance. Do NOT use for broad local SEO outside the profile, general content drafting, or nonlocal search strategy.
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
    - google_business_profile_optimization
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - google_business_profile_optimization_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Google Business Profile Expert

## Trigger Conditions

- User wants to optimize Google Business Profile
- User asks about Google Maps rankings
- User wants to improve local pack visibility
- User asks about GBP categories, services, photos, reviews, Q&A, or posts
- User wants to fix GBP completeness or relevance
- User asks how to beat local competitors in Maps
- User needs local prominence or relevance improvements
- User wants to connect GBP to local SEO strategy

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the GBP goal
- Determine whether the user wants:
    * better Maps rankings
    * better local pack visibility
    * stronger profile relevance
    * more calls / direction requests / clicks
    * better review performance
    * better category alignment
    * better service/product coverage
- If unclear, ask one question:
    "Are you trying to improve Maps rankings, profile completeness, or review performance?"

### STEP 2 - Check profile fundamentals
Analyze:
  - business name consistency
  - primary category
  - secondary categories
  - business description
  - website URL
  - phone number
  - address / service area
  - hours
  - attributes
  - opening date
  - business type
  - booking or contact links

### STEP 3 - Evaluate relevance signals
Check whether the profile includes:
  - category alignment with core services
  - services/products fully filled out
  - business description using relevant terms naturally
  - photos that match the business and category
  - Q&A coverage for common intent questions
  - posts that reinforce services or offers
  - review language that supports topic relevance

### STEP 4 - Evaluate prominence signals
Analyze:
  - review count
  - review velocity
  - average rating
  - owner responses
  - photo freshness
  - posting cadence
  - local brand mentions
  - citations and local backlinks that support the profile

### STEP 5 - Evaluate completeness signals
Check for missing or weak:
  - business description
  - services
  - products
  - photos
  - Q&A
  - business hours
  - special hours
  - booking links
  - attributes
  - profile updates

### STEP 6 - Identify competitor gaps
Compare against nearby/local competitors to find:
  - categories they use that you do not
  - services they list that you do not
  - review volume advantage
  - photo frequency advantage
  - post frequency advantage
  - Q&A advantage
  - stronger local authority signals

### STEP 7 - Prioritize fixes
Rank recommendations in this order:

  HIGH PRIORITY:
    - wrong or weak primary category
    - missing services/products
    - incomplete business description
    - poor review volume or response strategy
    - missing photos or stale photos
    - no Q&A coverage for core offerings

  MEDIUM PRIORITY:
    - posting cadence improvements
    - secondary category refinement
    - attribute optimization
    - booking/contact link cleanup
    - photo naming and organization

  LOWER PRIORITY:
    - minor wording edits
    - cosmetic profile improvements
    - optional attribute additions

### STEP 8 - Build specific recommendations
For each fix, include:
  - what is wrong or missing
  - why it affects Maps or local pack visibility
  - exactly where to update it in GBP
  - how to word or structure it
  - expected impact

### STEP 9 - Validate before returning
Check that:
  - recommendations are specific to Google Business Profile
  - no generic local SEO advice replaces GBP-specific guidance
  - prominence, relevance, and completeness are all addressed
  - suggestions are tied to Maps performance outcomes
  - competitor comparison is included if relevant

### STEP 10 - Build output
Present as:
  1. GBP visibility snapshot
     (relevance, prominence, completeness)
  2. Highest-impact fixes
     (what to change first)
  3. Category and services audit
     (alignment issues and missing services)
  4. Photo, review, and post strategy
     (how to strengthen prominence)
  5. Q&A and description improvements
     (how to improve relevance)
  6. ONE fastest GBP win
     (highest ROI action this week)
  7. ONE follow-up question

### GUARDRAILS:
- Never confuse GBP optimization with generic local SEO
- Never ignore categories
- Never ignore services/products
- Never ignore review strategy
- Never recommend changes without linking them to Maps visibility or profile relevance
- Always think in terms of relevance, prominence, and completeness

## Purpose

Provide procedural guidance to review Google Business Profile completeness and local visibility.

## When to use

- Use when the authorized intent is `site_audit` and the request is to review Google Business Profile completeness and local visibility.

## When NOT to use

- Do not use when the request belongs to `local_seo_optimization`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `google_business_profile_optimization_report`

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
