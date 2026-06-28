---
name: content-seo-settings
description: |
  Generates SEO metadata, slugs, Open Graph fields, FAQs, schema, and internal-link suggestions for a known content piece. Use when content exists or its purpose is defined. Do NOT use for drafting the article, sitewide technical audits, or unrelated configuration.
---

# Content SEO Settings Expert

## Trigger Conditions

- User wants SEO metadata for a content piece
- User asks for title tags and meta descriptions
- User wants FAQ schema for an article
- User wants Open Graph tags
- User asks for slug optimization
- User wants structured data for a blog post
- User asks for internal linking suggestions
- User wants to optimize all SEO settings after writing
- User asks how to configure SEO for a specific page

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the content piece
- Identify what has been written or needs SEO settings:
    * blog post
    * how-to guide
    * landing page
    * product page
    * comparison page
    * FAQ page
    * pillar page
- Identify the primary keyword and topic
- If unclear, ask one clarifying question:
    "What page or content piece needs SEO settings configured?"

### STEP 2 - Gather source context
- If a URL is provided, inspect the page content
- If content text is provided, use that as the source
- Use quickContext for brand/domain context
- Identify:
    * primary keyword
    * secondary keywords and entities
    * target audience
    * content format and structure
    * key questions the content answers

### STEP 3 - Generate title tag
- Place primary keyword near the front
- Keep it concise and specific to the page
- Make it compelling for clicks
- Differentiate from other pages on the site
- Avoid keyword stuffing
- Avoid generic or vague titles

### STEP 4 - Generate meta description
- Summarize the page value clearly
- Include primary keyword naturally
- Use active language and benefit framing
- Match the search intent
- Make it persuasive without being misleading
- Keep it concise and scannable

### STEP 5 - Generate URL slug
- Use the primary keyword as the base
- Keep it short and readable
- Use hyphens between words
- Remove stop words when possible
- Avoid dates unless content is time-sensitive
- Make it descriptive and permanent

### STEP 6 - Generate FAQ schema
- Identify the key questions the content answers
- Write clean question and answer pairs
- Keep answers concise and direct
- Use answer-ready format for AI citation potential
- Output valid JSON-LD FAQ schema
- Include only questions supported by the content

### STEP 7 - Generate Article schema
- Include:
    * headline
    * author
    * datePublished
    * dateModified
    * publisher
    * image
    * description
- Output valid JSON-LD Article schema
- Keep values aligned with visible page content

### STEP 8 - Generate Open Graph tags
- Include:
    * og:title
    * og:description
    * og:image
    * og:url
    * og:type
    * og:site_name
- Optimize for social sharing and preview appearance
- Keep og:title and og:description distinct from
  meta title and description when appropriate

### STEP 9 - Generate internal linking suggestions
- Identify related pages on the site that should
  link TO this content
- Identify pages this content should link OUT to
- Suggest anchor text for each internal link
- Prioritize links that:
    * support the pillar and cluster architecture
    * pass authority to important pages
    * improve topical relevance signals

### STEP 10 - Validate mentally before returning
Check for:
  - keyword in title, meta, slug, and schema
  - no duplicate metadata across pages
  - valid JSON-LD syntax
  - accurate schema values
  - internal links that make sense contextually
  - Open Graph tags complete and accurate

### STEP 11 - Build output
Present as:
  1. Title tag
  2. Meta description
  3. URL slug
  4. FAQ schema (JSON-LD)
  5. Article schema (JSON-LD)
  6. Open Graph tags
  7. Internal linking suggestions
     (pages to link from and to, with anchor text)
  8. ONE follow-up question if needed

### GUARDRAILS:
- Never generate metadata without reading the content first
- Never keyword-stuff titles, descriptions, or schema
- Never invent FAQ questions not supported by the content
- Never use the same slug pattern for every page
- Never skip internal linking suggestions
- Always keep schema values aligned with visible content
- Always generate all settings together as a complete package
- Always treat this as the final SEO layer after content is written
