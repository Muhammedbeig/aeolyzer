---
name: seo-search-intent
description: Guidelines for SEO search intent alignment, AEO optimization, evidence density, readability craft, UGC intelligence, and self-contained section structure. Use when writing or optimizing any content for search engines or AI answer systems. Do NOT use for technical crawling, analytics reporting, or writing unrelated to search.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - seo_planning
allowed_modes:
    - plan
    - read
capability_tags:
    - seo_search_intent
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - seo_search_intent_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# SEO + Search Intent

## SEO + SEARCH INTENT

- Identify the primary keyword or search query
- Align content type to intent: informational, commercial, navigational, comparison
- Check cannibalization risk against the user's existing content
- Cover adjacent questions that support topic authority
- Use question-based subheadings where useful
- Add concise, self-contained answers that can be extracted by search engines and AI systems
- Support claims with cited sources and clear language

## AEO (Answer Engine Optimization)

- AI engines (ChatGPT, Perplexity, Google AI Overviews) extract individual sections from articles.
  The best-performing sections lead with a concise 40-60 word direct answer, a self-contained
  statement that makes sense if quoted in isolation. Pattern: Fact, then interpretation, then implication.
- Sections that start with "Let's explore..." or reference other sections ("As mentioned above...")
  perform poorly in extraction. Self-contained sections get cited.
- Pages with paragraph-length summaries or key takeaways have ~35% higher AI snippet inclusion.
  For longer pieces, a "Key Takeaways" section (5-7 bullets) near the top or end can help.

## EVIDENCE DENSITY

- Content with specific statistics gets ~40% higher citation rates from AI engines. The "fluency +
  statistics" combination outperforms any single content strategy by 5.5%+.
- A good benchmark: at least one specific data point per major section, and 3-5 cited sources per
  1,000 words. "According to [X]...", "[Source] found that..." reads naturally.
- Replace vague quantifiers ("many companies", "significant improvement") with specific numbers
  when your research provides them.

## READABILITY CRAFT

- Shorter paragraphs (~40-60 words, roughly 2-3 sentences) improve both human scanability and AI
  extraction quality. Longer paragraphs create "muddy" embeddings that reduce retrieval confidence.
- Sentence rhythm matters: alternate short declarative statements (5-8 words) with medium explanatory
  ones (12-18 words). Monotonous sentence length makes prose feel flat.
- Question-phrased headings ("What Does X Cost?" vs "X Pricing") map to how users query AI chat
  and appear in People Also Ask boxes. Useful for informational and comparison sections.

## UGC INTELLIGENCE

- Real user language from Reddit, forums, and Quora reveals hidden intent better than keyword tools.
  A search like "site:reddit.com [topic]" surfaces the actual words people use, their pain points,
  and questions that existing articles don't answer. When you find good UGC language during research,
  echoing that phrasing in the article makes it resonate more authentically.

## SELF-CONTAINED SECTIONS

- Each H2 section should pass the "extraction test": readable in isolation without needing context
  from surrounding sections. Pronouns that reference unstated antecedents ("This approach...") and
  cross-references without context ("as discussed above") weaken extractability. A brief restatement
  is better than a dangling reference.

## Purpose

Provide procedural guidance to classify search intent and explain content-format implications.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to classify search intent and explain content-format implications.

## When NOT to use

- Do not use when the request belongs to `keyword_research`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `seo_search_intent_report`

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
