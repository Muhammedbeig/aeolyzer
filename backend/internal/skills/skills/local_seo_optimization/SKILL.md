---
name: local-seo-optimization
description: Plans local search improvements across Maps, business profiles, citations, NAP consistency, reviews, local pages, schema, and geo-targeted keywords. Use when the user wants visibility in a city or service area. Do NOT use for nonlocal SEO, generic content strategy, or unsupported location claims.
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
    - local_seo_optimization
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - local_seo_optimization_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Local SEO Optimizer

## Trigger Conditions

- User wants to rank in Google Maps or the local pack
- User asks about local rankings or "near me" visibility
- User wants to optimize a Google Business Profile
- User asks about citations, local directories, or NAP consistency
- User wants to improve visibility in a city, region, or service area
- User asks how to win local search against nearby competitors
- User wants geo-targeted keyword and page strategy
- User needs local landing page or service area optimization

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the local objective
- Determine whether the user wants:
    * Google Maps / local pack rankings
    * service area visibility
    * city-page rankings
    * review growth
    * citation consistency
    * local landing page optimization
    * multi-location strategy
- If unclear, ask one question:
    "Are you trying to rank in Maps, improve city-page visibility, or fix local citations?"

### STEP 2 - Establish business location context
- Confirm the business location, service area, or target city from project context or memory
- Identify whether the business is:
    * storefront
    * service-area business
    * multi-location brand
    * hybrid model
- Determine if the website structure matches the local SEO model

### STEP 3 - Audit local presence
Analyze:
  - Google Business Profile completeness
  - NAP consistency across site and citations
  - city/service area pages
  - location page structure
  - review volume and sentiment
  - local schema markup
  - local backlinks
  - embedded map usage
  - local content relevance

### STEP 4 - Evaluate local keyword targeting
Check whether the site targets:
  - city + service keywords
  - "near me" intent
  - neighborhood modifiers
  - service area modifiers
  - geo-specific long-tail queries
Flag if:
  - keywords are too broad
  - local pages are missing
  - page titles do not include location modifiers
  - content does not match local intent

### STEP 5 - Check Google Business Profile optimization
Evaluate:
  - business categories
  - business description
  - services/products
  - photos
  - Q&A
  - posts
  - reviews
  - UTM tracking
  - primary category relevance

### STEP 6 - Examine citations and directory footprint
Identify:
  - major citation inconsistencies
  - missing local directories
  - duplicate listings
  - incorrect address/phone formatting
  - weak category alignment in directories
  - opportunities in authoritative local sources

### STEP 7 - Analyze local page structure
For location/service pages, check:
  - unique location content
  - city-specific proof
  - local testimonials
  - service area mention
  - directions or map embed
  - local FAQs
  - schema markup
  - internal links from homepage and service pages

### STEP 8 - Prioritize by local impact
Rank fixes in this order:

  HIGH PRIORITY:
    - Google Business Profile issues
    - NAP inconsistencies
    - missing or weak city pages
    - no local schema
    - poor category targeting
    - missing reviews or weak review velocity

  MEDIUM PRIORITY:
    - local backlinks
    - service area page improvements
    - local FAQ expansion
    - local internal linking
    - photo/post optimization

  LOWER PRIORITY:
    - minor citation cleanup
    - directory polishing
    - small copy tweaks
    - map embed styling changes

### STEP 9 - Build specific recommendations
For each recommendation, include:
  - what to fix
  - why it matters locally
  - where it should happen
  - expected outcome
  - whether it affects Maps, organic local pages, or both

### STEP 10 - Validate before returning
Check that:
  - recommendations are tied to local ranking factors
  - location intent is clear
  - GBP, citations, and local pages are all considered
  - suggestions are specific and actionable
  - no generic SEO advice is mixed in without local relevance

### STEP 11 - Build output
Present as:
  1. Local visibility snapshot
     (GBP, citations, pages, reviews, schema)
  2. Highest-impact local fixes
     (what to do first)
  3. Citation and NAP issues
     (inconsistencies, missing listings, duplicates)
  4. Local page improvements
     (city pages, service area pages, FAQs, schema)
  5. Google Business Profile actions
     (categories, posts, photos, reviews, Q&A)
  6. ONE fastest local win
     (highest ROI action this week)
  7. ONE follow-up question

### GUARDRAILS:
- Never treat local SEO as generic SEO
- Never ignore Google Business Profile
- Never ignore NAP consistency
- Never recommend city pages without checking local intent
- Never skip citation and directory analysis
- Always consider Maps and organic local pages separately
- Always prioritize business location relevance over broad keyword volume

## Purpose

Provide procedural guidance to review local SEO consistency, relevance, and location signals.

## When to use

- Use when the authorized intent is `site_audit` and the request is to review local SEO consistency, relevance, and location signals.

## When NOT to use

- Do not use when the request belongs to `google_business_profile_optimization`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `local_seo_optimization_report`

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
