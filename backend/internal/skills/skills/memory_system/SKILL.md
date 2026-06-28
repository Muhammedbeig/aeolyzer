---
name: memory-system
description: Use to read and update persistent brand memory and tone documents across content sessions. Covers when to call readMemoryDoc and proposeMemoryUpdate, how to identify patterns worth saving, and when NOT to propose updates. Trigger when a user shares a new preference, corrects tone or style repeatedly, or after any significant style or structural rework — check the tone doc and propose missing rules before moving on. Do NOT use for silent memory writes, one-off preferences, or content drafting.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - memory_tone_management
allowed_modes:
    - plan
    - read
capability_tags:
    - memory_system
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - memory_system_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Memory System

You have two persistent markdown documents injected into your context above:

- `<memory>`: Brand identity, strategic priorities, audience, competitors, content history, learned preferences
- `<tone>`: Writing voice, formatting rules, word preferences, banned phrases, style guidelines

These documents are the source of truth for the user's preferences. They persist across sessions.

## Tools

- **readMemoryDoc(slug):** Re-read the latest version of a memory doc (useful after accepting a proposal)
- **proposeMemoryUpdate(slug, section, proposedContent, reason):** Propose a change. Renders an accept/reject card in chat. NEVER update memory silently. Always propose and let the user decide.

## When to propose updates

- User shares a new preference ("always use Oxford commas", "our audience is enterprise CTOs")
- You notice a pattern in their feedback that should be codified
- After a writing session where the user corrected your tone or style repeatedly
- When the user tells you something about their brand, competitors, or strategy
- After you rewrite or restyle content: if the user asked you to match their tone, remove bold formatting, shorten paragraphs, etc., and the changes reveal style rules not already captured in the tone doc, propose adding those rules. Example: you rewrote a piece to be "more direct, no bold, shorter paragraphs" -> propose a tone update capturing those specific preferences so future content gets it right the first time.
- When you apply a correction more than once in the same session, that's a pattern worth saving.

Be proactive: if you just did significant style or structural work based on the user's feedback, check the tone doc (readMemoryDoc) and propose any missing rules BEFORE moving on. Don't wait for the user to ask.

## When NOT to propose

- For one-off instructions that only apply to the current piece
- For temporary preferences ("make this one more casual" doesn't mean all content should be casual)
- If the information is already in the memory/tone doc
- If the preference is too vague to be actionable (e.g. "make it better"), ask for specifics first

## Purpose

Provide procedural guidance to apply approved tone summaries and propose memory changes for review.

## When to use

- Use when the authorized intent is `memory_tone_management` and the request is to apply approved tone summaries and propose memory changes for review.

## When NOT to use

- Do not use when the request belongs to `editorial_voice`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `memory_system_report`

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
