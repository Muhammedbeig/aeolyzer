---
name: outline-structure
description: |
  Guidance for heading hierarchy, section flow, and formatting variety when writing articles. Use when writing, planning, or structuring multi-section articles, blog posts, or long-form content — even if the user doesn't say "outline" or "structure" explicitly. Do NOT use for short-form copy, metadata-only tasks, or final publishing.
---

# Outline Structure

## HEADING HIERARCHY

- The title field (set via generateTitle or the canvas title editor) IS the page's H1.
- Article body content MUST start at H2. Never write a markdown # H1 heading inside the body.
- Use ## H2 for major sections, ### H3 for subsections within those.
- If the plan includes a section titled "Introduction", it does NOT get an H2; the intro flows
  directly as prose under the title. The first H2 should be the first named section after the intro.

## DUPLICATE H2 BUG PREVENTION (critical)

- When you call writeSection with a sectionTitle, the tool AUTOMATICALLY prepends "## sectionTitle"
  to your markdown. Do NOT also start your markdown with "## ..." or you will get a DUPLICATE H2.
- Rule: If you pass sectionTitle, your markdown must NOT begin with any ## heading.
  Start the markdown with the body text, H3 subheadings, or other content.
- If you omit sectionTitle, then you are responsible for including the ## H2 in your markdown.

## SECTION FLOW

1. Review the full plan (if it exists). Note per-section word allocations, key points,
   evidence, and hiddenIntent.
2. Start with the introduction/hook - establish the thesis and why this matters NOW.
   If the plan has a golden thread, weave it into the opening.
3. Write each body section. Before each one, announce it in chat with the planned
   word count and key points.
4. BEFORE each section, check: "How many words have I written so far vs the target?"
5. End with a conclusion that ADVANCES the argument (not just summarizes it).
   What should the reader DO next?
6. After the last section, report the actual total from summing writeSection wordCounts.
7. The user can read and react to earlier sections while you write later ones.

## FORMATTING VARIETY, REQUIRED (THIS IS THE #1 QUALITY SIGNAL)

Text walls are the single most common failure mode. Readers skim. AI engines extract
structured elements. An article that's all paragraphs is an article that fails both
audiences. Treat this section as law.

### HARD RULES

- NEVER write more than 2 consecutive plain paragraphs without a structural break.
  After 2 paragraphs, you MUST insert one of: bullet list, numbered list, table, blockquote,
  bold callout, H3 subheading, or key takeaway box. No exceptions.
- Every H2 section MUST contain at least TWO non-paragraph elements (not just one).
  Example: an H3 subheading + a bullet list. Or a table + a blockquote. Variety matters.
- Aim for a structural element every 100-150 words. If you've written 150+ words of pure
  prose, you've gone too long without a break.

### H3 SUBHEADINGS: USE THEM AGGRESSIVELY

- Any H2 section longer than 200 words SHOULD have at least one H3 to break it up.
- H3s aren't just cosmetic. They create extraction points for AI engines, anchor links
  for navigation, and visual breathing room for readers.
- Pattern: H2 intro paragraph (2-3 sentences) -> H3 first subtopic -> H3 second subtopic.
  This is almost always better than H2 -> 6 paragraphs.

### STRUCTURAL ELEMENTS AND WHEN TO USE THEM

- Bullet lists: 3+ parallel items, feature lists, key points, requirements, takeaways
- Numbered lists: Sequential steps, ranked items, prioritized recommendations
- Tables: Side-by-side comparisons (2+ items with 2+ dimensions). Use markdown pipe tables.
- Blockquotes: Key quotes with attribution, critical definitions, important callouts
- Bold callouts: Single-line emphasis for key stats, surprising facts, or critical warnings.
  Pattern: **Key insight:** [the insight]. Stands out from surrounding prose.
- Key takeaway boxes: At the start or end of major sections. Use a blockquote with
  "> **Key takeaway:** ..." or "> **TL;DR:** ..." formatting.

### SECTION PACING

- Opening section (intro): Hook paragraph -> bold callout or key stat -> context paragraph
  -> short list of what the article covers. NEVER open with 4+ paragraphs of pure prose.
- Body sections: H2 -> 1-2 paragraph intro -> H3 or structural element -> content ->
  H3 or structural element -> content. Each section should feel like it has internal architecture.
- Closing section: Short. 1-2 paragraphs max. A bold callout or key takeaway, then the CTA.
