---
name: post-write-checklist
description: |
  A 12-point quality checklist to run after completing all article sections. Covers structure, text walls, H3 usage, link counts, redundancy, completeness, tables, quotes, formatting variety, em dashes, duplicate headings, and SEO suggestions. Use immediately after finishing any article or long-form content piece. Do NOT use before the draft is complete, for content planning, or for research.
---

# Post-Write Checklist

## POST-WRITE QUALITY CHECKLIST

After finishing ALL writeSection calls, run this checklist before reporting completion.
If multiple checks fail across many sections, use writeSection("replace") with the full revised article.
If only 1-2 isolated issues, use proposeEdit for targeted fixes.

1. STRUCTURE: Every section from the plan must produce a corresponding H2 in the output.
   Count your H2s against the plan's section list. Missing H2 = missing section = not done.
2. TEXT WALL SCAN (ZERO TOLERANCE):
   Scan the ENTIRE article top to bottom. If you find ANY stretch of 3+ consecutive plain
   paragraphs without a list, table, H3, blockquote, or bold callout between them, it FAILS.
   Fix it immediately. The opening section MUST include a structural break within the first
   150 words. If your intro is 4 paragraphs of pure prose, rewrite it with a bold callout
   or key stat callout after the hook.
3. H3 CHECK: Any H2 section longer than 250 words MUST contain at least one H3 subheading.
   If a section is 400+ words with zero H3s, it's a wall. Break it up.
4. LINK COUNT (HARD MINIMUM):
   Count the total number of inline hyperlinks [text](URL) in the article.
   - Under 1,000 words: minimum 3 links
   - 1,000-2,000 words: minimum 5 links
   - 2,000-3,000 words: minimum 8 links
   - 3,000+ words: minimum 10 links
   An article with fewer links than the minimum is NOT finished. Go back and add source
   citations, reference links, and contextual links from your research.
   An article with ZERO links is an automatic failure regardless of word count.
5. REDUNDANCY: The core thesis is stated ONCE in the introduction. If the same point appears
   in 2+ paragraphs (even with different wording), cut the duplicate. Each paragraph must
   advance the argument, not restate it.
6. COMPLETENESS: Compare your running word total to the target. If you're more than 20% short,
   sections are missing. Go back and write them.
7. TABLES: If you wrote a markdown table, verify it has proper pipe-delimited rows with a
   header separator line (|---|---|). If the table rendered incorrectly, replace it with a
   comparison list (bold label + description pairs).
8. QUOTES: Blockquotes with named attribution MUST include a source link. Never fabricate
   quotes. If you don't have a direct quote URL, use inline paraphrase instead of a blockquote
   with a fake attribution.
9. STRUCTURAL VARIETY: No two consecutive H2 sections should use the same primary format.
   If section 2's main element was a bullet list, section 3 should lead with a table, blockquote,
   numbered list, or comparison. Monotonous formatting is almost as bad as no formatting.
10. EM DASH SCAN: Search the entire output for em dashes (—). If you find ANY, replace them
    with commas, semicolons, colons, parentheses, or split into separate sentences. Zero tolerance.
11. DUPLICATE H2 SCAN: Check that no section has two consecutive H2 headings. If writeSection
    was called with sectionTitle AND the markdown also started with ##, you have a duplicate. Fix it.
12. SEO SUGGESTION: After all checks pass, suggest SEO optimization: "Content's done. Want me to
    set up SEO settings? I'll generate meta title, description, FAQs, and schema markup."
    This is a natural next step; most users want it but won't think to ask.
