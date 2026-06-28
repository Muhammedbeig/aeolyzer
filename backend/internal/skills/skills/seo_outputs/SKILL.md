---
name: seo-outputs
description: |
  Generates all SEO output fields after article writing: meta title, meta description, slug, OG tags, FAQs, schema markup, internal link suggestions, and social variants. Use after completing any article to produce the full SEO package. Do NOT use before article content exists, for full drafting, or for technical site audits.
---

# SEO Outputs

## SEO OUTPUTS AFTER WRITING

### meta title

- clear, searchable, strong CTR
- include primary keyword near front
- max 60 characters
- format: [Primary Keyword] - [Value Prop or Context]
- never overwrite if already set

### meta description

- concise summary with click intent
- include primary keyword naturally
- max 155 characters
- answer: what will the reader get from this page?
- never overwrite if already set

### slug

- short, clean, lowercase, hyphenated
- no stop words (a, the, and, of, in)
- matches primary keyword
- example: /premier-league-top-scorers-2026
- never overwrite if already set

### og title

- social-friendly version of meta title
- can be slightly longer or more engaging
- optimized for shares and clicks on social

### og description

- social-friendly summary
- 1-2 sentences
- designed for social feed context, not search

### faqs

- 3-5 real questions tied to article intent
- questions must match what readers actually search
- answers must be short, direct, and self-contained
- no fluff, no filler
- format: Q + A pairs
- each answer must be 1-3 sentences max

### schema markup

- choose schema type based on page content:
    Article      : general informational articles
    BlogPosting  : blog posts and opinion pieces
    NewsArticle  : match reports, breaking news, transfers
    FAQPage      : pages with question and answer content
    WebPage      : general landing or topic pages
    SportsEvent  : match previews or event pages
- build valid JSON-LD using:
    article title
    meta description
    page URL
    publisher name and domain
    date published / date modified
    FAQ entities if FAQs are generated
- output format: JSON-LD
- injection: <script type="application/ld+json">

### internal link suggestions

- find related pages on the same site (livesoccer24.com)
- match by topic, entity, or keyword overlap
- suggest anchor text for each link
- connect topically related pages only
- output format:
    page_title  : string
    url         : string
    anchor_text : string
    relevance   : why this link fits contextually

### social title and description variants

- twitter/x card title
- twitter/x card description
- optimized for social sharing context

## VALIDATION RULES

- do not overwrite any existing seo field values
- meta title must align with article title and intent
- meta description must match article content
- slug must be clean, short, and keyword-aligned
- faqs must be answerable and non-duplicate
- schema type must match content type
- schema JSON-LD must be valid and complete
- internal links must be relevant and safe
- no competitor domains in any links
- noIndex and noFollow must be false unless explicitly set
