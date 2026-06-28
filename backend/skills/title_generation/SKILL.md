---
name: title-generation
description: |
  Generates specific, differentiated titles by combining surface, emotional, and business intent. Use when a topic and angle are known and a content title is needed. Do NOT use for outlines, full drafts, metadata, or unrelated naming tasks.
---

# Title Generation

The best-performing content doesn't just match what people typed - it matches what people FEEL but didn't type. Every search query has up to three intent layers. Your job is to identify all three and resolve them in the content.

## THE THREE INTENT LAYERS

1. **SURFACE INTENT** - what the query literally asks.
   Example: "paypal international fees" -> "What are PayPal's international fees?"
   This is what every competitor already covers. Matching it is table stakes.

2. **EMOTIONAL INTENT** - the anxiety, desire, frustration, or need state driving the search.
   Example: "paypal international fees" -> "I hate paying these fees and want to avoid them"
   This is what the searcher FEELS but doesn't type. It's the unstated pain point, aspiration, or fear behind the query. Most articles miss this entirely.

3. **BUSINESS INTENT** - the ROI or justification the reader needs (especially for B2B/professional content).
   Example: "optimize for ai overviews" -> "Is this worth my team's time and budget?"
   Professional searchers almost always need to justify the time investment. If the surface intent is "how to do X", there's a sublayer of "is X worth doing?"

Not every query has all three layers. A casual how-to may only have surface + emotional. B2B and professional content almost always has all three.

## THE GOLDEN THREAD

The golden thread is a single narrative line that resolves all identified layers. It runs from the title through the introduction, through the body, to the conclusion. It's what makes the content feel cohesive rather than like a collection of H2s.

Example (adapt to the user's industry): "The cost panic is real, but the margin story tells the opposite. Early adopters are seeing 3-5x higher ROI than the old approach."
- Resolves surface intent (what's happening with [trend])
- Resolves emotional intent (am I falling behind? -> no, you're gaining)
- Resolves business intent (is this worth investing in? -> yes, here's the ROI)

## TITLE GENERATION - HIDDEN INTENT IN THE TITLE

The title is where hidden intent has the most impact on performance. Pattern:

**[Surface intent answer] + [Emotional intent resolution]**

Canonical example:
- **BAD:** "What are PayPal's International Fees" (surface only - matches every competitor)
- **GOOD:** "What are PayPal's International Fees and How to Avoid Them" (surface + emotional)

The second half of the title is the hidden intent hook. It's what drives higher CTR because it speaks to what the searcher actually wants, not just what they typed.

When proposing titles (in `proposePlan` or `generateTitle`):
- Always consider: does this title address only the surface, or does it also capture the emotional/hidden need?
- If you identified a hidden intent, the title MUST reflect it.
- The hidden intent hook typically comes after a conjunction: "and How to...", "and Why...", "and What to Do About It", "Without...", etc.
- Not every title needs this pattern. Short-form content (LinkedIn, YouTube descriptions) may use the hidden intent as the entire angle rather than appending it.

## TITLE RULES

- Call `generateTitle` as SOON as you know the topic, ideally right after the user confirms what they want to write about, BEFORE writing the first section.
- The title shows in the header bar and sidebar; leaving it as "Untitled Article" makes the workspace feel empty and unprofessional.
- In Plan mode, `proposePlan` already sets the title. But in Write mode, you must call `generateTitle` explicitly once the topic is clear.
- You can always update the title later with another `generateTitle` call if the angle evolves.
