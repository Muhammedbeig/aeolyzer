---
name: memory-system
description: |
  Use to read and update persistent brand memory and tone documents across content sessions. Covers when to call readMemoryDoc and proposeMemoryUpdate, how to identify patterns worth saving, and when NOT to propose updates. Trigger when a user shares a new preference, corrects tone or style repeatedly, or after any significant style or structural rework — check the tone doc and propose missing rules before moving on. Do NOT use for silent memory writes, one-off preferences, or content drafting.
---

# Memory System

You have two persistent markdown documents injected into your context above:

- `<memory>`: Brand identity, strategic priorities, audience, competitors, content history, learned preferences
- `<tone>`: Writing voice, formatting rules, word preferences, banned phrases, style guidelines

These documents are the source of truth for the user's preferences. They persist across sessions.

## Tools

- **readMemoryDoc(slug):** Re-read the latest version of a memory doc (useful after accepting a proposal)
- **proposeMemoryUpdate(slug, section, proposedContent, reason):** Propose a change. Renders an accept/reject card in chat. NEVER update memory silently. Always propose and let the user decide.

## When to propose updates

- User shares a new preference ("always use Oxford commas", "our audience is enterprise CTOs")
- You notice a pattern in their feedback that should be codified
- After a writing session where the user corrected your tone or style repeatedly
- When the user tells you something about their brand, competitors, or strategy
- After you rewrite or restyle content: if the user asked you to match their tone, remove bold formatting, shorten paragraphs, etc., and the changes reveal style rules not already captured in the tone doc, propose adding those rules. Example: you rewrote a piece to be "more direct, no bold, shorter paragraphs" -> propose a tone update capturing those specific preferences so future content gets it right the first time.
- When you apply a correction more than once in the same session, that's a pattern worth saving.

Be proactive: if you just did significant style or structural work based on the user's feedback, check the tone doc (readMemoryDoc) and propose any missing rules BEFORE moving on. Don't wait for the user to ask.

## When NOT to propose

- For one-off instructions that only apply to the current piece
- For temporary preferences ("make this one more casual" doesn't mean all content should be casual)
- If the information is already in the memory/tone doc
- If the preference is too vague to be actionable (e.g. "make it better"), ask for specifics first
