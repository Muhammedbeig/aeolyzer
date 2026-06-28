---
name: brief-building
description: |
  Builds a content brief covering topic, angle, audience, intent, CTA, keywords, length, subtype, and notes. Use before drafting begins or when planning article requirements. Do NOT use for writing sections, researching sources, or publishing content.
---

# Brief Building

Before writing begins, you need a brief. The brief is the contract between you and the user - it defines what you're writing, for whom, and why. Without it, you're guessing.

## MINIMUM FIELDS (must have before writing)

- Topic: what are we writing about?
- Angle: what is the unique perspective or thesis?
- Audience: who is reading this?
- Intent: what should the reader know, feel, or do after reading?
- CTA: where should the reader go next?

## IDEAL FIELDS (collect when possible)

- Keywords: target search terms or search query
- Length: target word count based on content type and goal
- Subtype: opinion | comparison | topic-guide | news-analysis | how-to | alternatives
- Notes: anything specific the user wants included or avoided

## SUBTYPE SELECTION

- opinion: user has a strong POV, wants to take a position
- comparison: evaluating two or more options side by side
- topic-guide: comprehensive coverage of a subject
- news-analysis: breaking down a recent event or development
- how-to: step-by-step instructional content
- alternatives: positioning options against a category leader

## RULES

- Call `updateBrief()` before the first `writeSection()` call - always
- If the user skips brief questions and says "just write it", infer the brief from context and save it anyway
- If topic is clear but angle is missing, suggest 2-3 angle options before writing
- Never start writing without at least: topic, audience, and intent
- The brief is saved and visible alongside the canvas - treat it as a live document
- Update the brief if the user changes direction mid-article
