---
name: schema-generation
description: |
  Generates page-appropriate JSON-LD from verified content and business facts. Use when the user requests schema markup, rich-result support, or structured data for a known page. Do NOT use when required facts are missing, for crawler directives, or for unsupported schema types.
---

# Schema Markup Expert

## Trigger Conditions

- User asks for schema markup
- User asks for structured data
- User wants rich snippets
- User wants to help AI understand a page
- User asks for JSON-LD
- User wants FAQ, HowTo, Article, Product, Organization, or LocalBusiness schema
- User asks how to improve eligibility for search features
- User wants content made more machine-readable

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the page or entity type
- Determine what the user wants to mark up:
    * article
    * blog post
    * FAQ
    * how-to
    * product
    * organization
    * local business
    * service
    * breadcrumb
    * video
    * event
    * course
    * recipe
- If unclear, ask one clarifying question:
    "What type of page or entity are you marking up?"

### STEP 2 - Gather the source material
- If the user provides a URL, inspect the page content first
- If the user provides text, use that text as the source
- If the user asks for site-wide schema, identify the main templates/pages
- Use quickContext when brand context matters

### STEP 3 - Choose the correct schema type
- Match the page purpose to the best schema.org type
- Do NOT force irrelevant schema types
- Prefer the simplest valid schema that describes the page accurately
- If multiple types apply, combine them only when they are naturally connected

### STEP 4 - Build structured data
- Output valid JSON-LD
- Include required properties for the schema type
- Include recommended properties when available
- Keep values consistent with the page content
- Use absolute URLs where appropriate
- Ensure names, headlines, authors, dates, and descriptions match the page

### STEP 5 - Optimize for search and AI systems
- Make the markup explicit, clean, and unambiguous
- Use schema to support:
    * rich results
    * entity understanding
    * AI citation confidence
    * page classification
- For FAQ and HowTo, structure answers clearly
- For Organization, include brand identity details
- For Article, include headline, author, datePublished, dateModified, and publisher
- For Product, include name, image, description, offers, and aggregateRating when valid

### STEP 6 - Validate mentally before returning
Check for:
  - valid JSON syntax
  - correct type selection
  - property completeness
  - alignment with on-page content
  - no invented facts
  - no markup spam

### STEP 7 - Build output
Present as:
  1. Recommended schema type or types
  2. JSON-LD code block
  3. Short explanation of why it fits
  4. Optional notes on where to place it or how to test it
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never invent product ratings, reviews, or prices
- Never add FAQ questions not supported by the page
- Never use schema as a substitute for missing content
- Never overcomplicate simple pages
- Always keep markup aligned with visible page content
