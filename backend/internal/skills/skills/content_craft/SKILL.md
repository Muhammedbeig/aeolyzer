---
name: content-craft
description: Use when writing or editing long-form content that needs to perform well in traditional search and AI answer engines. Covers AEO (Answer Engine Optimization), revenue case framing, evidence density, readability craft, UGC intelligence, and self-contained section structure. Trigger for SEO articles, blog posts, B2B content, or any piece where performance and citation rates matter. Do NOT use for short coordination messages, metadata-only tasks, or raw research collection.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - draft_article
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - content_craft
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - content_craft_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Content Craft

You have deep knowledge of what makes long-form content perform well in both traditional search and AI answer engines. Use this knowledge with editorial judgment; apply what fits the piece, skip what doesn't. Not every article needs every technique.

## THINGS YOU KNOW ABOUT HIGH-PERFORMING CONTENT

### AEO (Answer Engine Optimization)

- AI engines (ChatGPT, Perplexity, Google AI Overviews) extract individual sections from articles. The best-performing sections lead with a concise 40-60 word direct answer, a self-contained statement that makes sense if quoted in isolation. Pattern: Fact, then interpretation, then implication.
- Sections that start with "Let's explore..." or reference other sections ("As mentioned above...") perform poorly in extraction. Self-contained sections get cited.
- Pages with paragraph-length summaries or key takeaways have ~35% higher AI snippet inclusion. For longer pieces, a "Key Takeaways" section (5-7 bullets) near the top or end can help.

### Revenue Case / Business Justification

- The most compelling B2B/professional content answers "why should my business care?" with hard data early in the piece, not buried at the end. Conversion rates, ROI case studies, market-size data, cost-of-inaction figures. This reframes the reader from anxiety to opportunity.
- When competitors lead with fear/loss framing, leading with opportunity data is a powerful differentiator.
- During research, searches like "[topic] conversion rate", "[topic] ROI case study", "[topic] revenue impact" often surface this data. Industry reports (Gartner, Forrester, HubSpot) are good sources.

### Evidence Density

- Content with specific statistics gets ~40% higher citation rates from AI engines. The "fluency + statistics" combination outperforms any single content strategy by 5.5%+.
- A good benchmark: at least one specific data point per major section, and 3-5 cited sources per 1,000 words. "According to [X]...", "[Source] found that..." reads naturally.
- Replace vague quantifiers ("many companies", "significant improvement") with specific numbers when your research provides them.

### Readability Craft

- Shorter paragraphs (~40-60 words, roughly 2-3 sentences) improve both human scanability and AI extraction quality. Longer paragraphs create "muddy" embeddings that reduce retrieval confidence.
- Sentence rhythm matters: alternate short declarative statements (5-8 words) with medium explanatory ones (12-18 words). Monotonous sentence length makes prose feel flat.
- Question-phrased headings ("What Does X Cost?" vs "X Pricing") map to how users query AI chat and appear in People Also Ask boxes. Useful for informational and comparison sections.

### UGC Intelligence

- Real user language from Reddit, forums, and Quora reveals hidden intent better than keyword tools. A search like "site:reddit.com [topic]" surfaces the actual words people use, their pain points, and questions that existing articles don't answer. When you find good UGC language during research, echoing that phrasing in the article makes it resonate more authentically.

### Self-Contained Sections

- Each H2 section should pass the "extraction test": readable in isolation without needing context from surrounding sections. Pronouns that reference unstated antecedents ("This approach...") and cross-references without context ("as discussed above") weaken extractability. A brief restatement is better than a dangling reference.

---

Apply these patterns when they serve the piece. A 3,000-word SEO article benefits from most of them. An 800-word blog post might use one or two. Trust your editorial judgment.

## Purpose

Provide procedural guidance to improve prose clarity, specificity, rhythm, and reader value.

## When to use

- Use when the authorized intent is `draft_article` and the request is to improve prose clarity, specificity, rhythm, and reader value.

## When NOT to use

- Do not use when the request belongs to `writing`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `content_craft_draft`
- `quality_summary`

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
