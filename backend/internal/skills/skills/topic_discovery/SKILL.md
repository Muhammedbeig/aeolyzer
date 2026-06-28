---
name: topic-discovery
description: Finds strategic content topics by examining brand context, competitor coverage, industry news, content gaps, and site overlap. Use when the user needs topic options or angles before a brief exists. Do NOT use for drafting, editing, or publishing content.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - topic_discovery
allowed_modes:
    - plan
    - read
capability_tags:
    - topic_discovery
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - topic_discovery_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# TOPIC DISCOVERY PROCESS (when user needs ideas)

Think like a content strategist being paid $300/hr. You wouldn't walk into a meeting and say "what about AI trends?" You'd arrive with a brief, having already researched the client's competitive landscape, recent news, and content gaps.

## Step 1: GATHER CONTEXT (do this silently, fast)

- Call `quickContext` to get: brand name, domain, industry, competitors, positioning
- This tells you WHO they are and WHO they compete with

## Step 2: RESEARCH THE LANDSCAPE (2-3 targeted searches)

Don't search "[industry] trending topics" - that gives generic garbage. Instead, search strategically:

- "[competitor name] blog 2026" (what are their competitors publishing?)
- "[industry] news 2026" (what just happened that's worth commenting on?)
- "site:[your-domain.com] blog" (what has the USER already written?)
  (so you don't suggest topics they've covered)

Look for GAPS: topics competitors cover that the user doesn't.
Look for HOOKS: recent news, launches, regulation changes, industry shifts.
Look for ANGLES: underserved perspectives on hot topics.

## Step 3: SYNTHESIZE INTO STRATEGIC SUGGESTIONS

Each suggestion must have THREE components:
1. THE TOPIC: specific, not vague
2. THE HOOK: why NOW (timeliness) or why THIS ANGLE (differentiation)
3. THE STRATEGIC VALUE: what publishing this does for the brand (positions them as X, captures audience searching for Y, counters competitor Z's narrative)

**BAD suggestions** (generic, zero strategic value):
- "The future of AI in [industry]"
- "Top 10 tips for 2026"
- "How to improve your strategy"

**GOOD suggestions** (specific, timely, strategically justified; adapt to the user's industry):
- "[Specific development]: what [target role] need to know before [deadline].
  Hook: [competitors X, Y] haven't covered the [specific angle] yet.
  This positions you as the authority in the space."
- "Why [metric the audience cares about] dropped after [recent event].
  Hook: [competitor] published a guide but missed the [differentiated angle] entirely. You'd be first to cover this."
- "[Competitor] just launched [feature]. Here's your counter-narrative.
  Hook: their announcement is getting traction, and a thoughtful response piece would capture the audience while framing your approach as superior."

## Step 4: PRESENT WITH CONFIDENCE

- Lead with your top recommendation and why: "I'd start with X - here's why."
- Present 2-3 options, each with the hook and strategic value.
- DON'T hedge. You're the strategist. Have a point of view.

## THE GAP MAP MINDSET

When suggesting topics, angles, or evaluating what to write, think in terms of GAPS:
- What are competitors publishing that the user ISN'T covering?
- What topics are getting traction in the industry with NO good coverage yet?
- What angle on a hot topic has everyone missed?
- What does the user's audience need to know that nobody is explaining well?

## FOUR TYPES OF CONTENT OPPORTUNITIES (prioritized)

1. COMPETITIVE GAPS: Competitor X published about [topic] but missed [angle].
   User can own that angle. Highest strategic value.
2. TIMING GAPS: Something just happened (launch, regulation, trend) and nobody in the user's space has published a thoughtful take yet. First-mover advantage.
3. DEPTH GAPS: Existing coverage is shallow (press release summaries, generic overviews).
   User can go deeper with original analysis, data, or expert perspective.
4. AUDIENCE GAPS: Content exists but it's written for the wrong audience.
   E.g. technical content exists but nothing for the decision-maker, or vice versa.

## NEVER

- Give topic suggestions without researching the competitive landscape first
- Suggest topics that sound like 2022 ChatGPT outputs ("Top 10 X", "The Future of Y")
- Present options without strategic justification (every option needs a "why")
- Ask questions with generic options when you have project/research context
- Say "Let me ask you a few questions" when they're already unsure
- Skip the research step when the user needs inspiration
- Re-ask about a topic the user already chose
- Rephrase a settled question with different wording

## Purpose

Provide procedural guidance to identify audience questions, content gaps, and defensible topic candidates.

## When to use

- Use when the authorized intent is `topic_discovery` and the request is to identify audience questions, content gaps, and defensible topic candidates.

## When NOT to use

- Do not use when the request belongs to `content_ideas`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `topic_discovery_report`

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
