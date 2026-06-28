---
name: sitemap-generation
description: |
  Generates or reviews an XML sitemap from verified canonical, indexable public URLs and sitemap limits. Use when the user requests sitemap creation, cleanup, splitting, or validation. Do NOT use for robots.txt, llms.txt, redirects, or inventing URLs.
---

# Sitemap Generator Expert

## Trigger Conditions

- User asks for a sitemap
- User asks how to submit URLs to search engines
- User wants to improve crawlability
- User asks about XML sitemap structure
- User wants sitemap index files for large sites
- User wants image, video, or news sitemaps
- User asks how to help search engines discover pages
- User asks how to keep sitemaps updated

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the site and sitemap goal
- Identify whether the user needs:
    * a standard XML sitemap
    * a sitemap index
    * image sitemap
    * video sitemap
    * news sitemap
    * a submission plan
- If unclear, ask one clarifying question:
    "What type of sitemap do you need?"

### STEP 2 - Gather page inventory context
- Use quickContext when available for site/domain context
- Determine which URLs should be included
- Prefer indexable, canonical, public URLs only
- Exclude duplicates, parameter URLs, private pages, and non-canonical pages
- For large sites, separate by section or template if needed

### STEP 3 - Assign priority and change frequency
- Homepage: priority 1.0, changefreq daily
- Key landing pages: priority 0.8-0.9
- Blog posts: priority 0.6-0.7, changefreq weekly
- Static pages: priority 0.5-0.8, changefreq monthly
- Category/archive pages: priority 0.4-0.6
- Low-value pages: priority 0.3-0.4
- Use lastmod values that reflect real content updates

### STEP 4 - Build sitemap structure
- Output valid XML
- Use proper sitemap namespace
- Include one <url> block per URL
- Add <loc>, <lastmod>, <changefreq>, and <priority>
- If site is large, provide sitemap index structure
- Keep file size under limits
- Make sure URLs are absolute and canonical

### STEP 5 - Optimize for discovery
- Ensure important pages are represented
- Group sitemaps logically if needed
- Include only URLs that should be indexed
- Add image/video extensions only where relevant
- Mention submission to Google Search Console and other search engines

### STEP 6 - Validate mentally before returning
Check for:
  - XML syntax correctness
  - canonical URL consistency
  - realistic priority values
  - accurate lastmod values
  - sitemap size and structure
  - no blocked or noindex URLs

### STEP 7 - Build output
Present as:
  1. XML sitemap content or sitemap index
  2. Brief explanation of the structure and why it fits
  3. Submission instructions
  4. Optional maintenance notes
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never include blocked, noindex, or canonicalized-away URLs
- Never invent lastmod dates
- Never overvalue low-priority pages
- Never omit sitemap index for very large sites
- Always keep the sitemap aligned with the site's actual indexable URLs
