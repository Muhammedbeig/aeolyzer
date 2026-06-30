# Frontend Agent Rules

> All design decisions (colors, spacing scale, typography, brand tokens) live in `design.md`.  
> This file governs **how** the agent works, not what things look like.

---

## 1. Tailwind Usage

- **Tailwind is the only styling system.** Never write inline `style={{}}` objects or separate CSS files for component-level styles. If a utility does not exist, extend `tailwind.config` — do not reach for a workaround.
- Compose classes in a consistent order: layout → sizing → spacing → typography → visual effects → state variants → responsive prefixes.
- Prefer semantic component extraction over chains of 15+ utility classes on one element. Extract a named component or use `@apply` in a component-scoped layer only when the class list becomes unreadable.
- Never hardcode arbitrary pixel values (e.g. `w-[347px]`) unless the value is genuinely one-off and cannot be expressed with the design scale. Document why when you do.
- Always use the project's design tokens (defined in `tailwind.config`) for spacing, radius, shadow, etc. Do not invent new arbitrary values if a token already covers the need.
- Before writing any frontend code use your this skill "frontend-patterns".
---

## 2. Visual Testing with Playwright — Non-Negotiable

Every UI change **must** be verified visually before the task is considered done. "It compiles" is not done. "It looks right in Playwright" is done.
- Before writing any frontend code use your this skill "frontend-patterns".

### 2.1 Required Test Coverage for Every Change

For each component or page touched, run Playwright tests that cover:

| Scenario | What to assert |
|---|---|
| Default / idle state | Screenshot matches baseline or looks correct on inspection |
| Hover state | Element changes appearance as expected (cursor, background, border, text) |
| Focus state | Focus ring is visible and accessible |
| Active / pressed state | Visual feedback is present |
| Disabled state | Element is visually muted and non-interactive |
| Loading / skeleton state | Placeholder renders without layout shift |
| Empty state | Empty UI renders gracefully, no broken layout |
| Error state | Error message or styling renders correctly |
| Responsive breakpoints | `mobile` (375 px), `tablet` (768 px), `desktop` (1280 px), `wide` (1536 px) |
| Dark mode (if supported) | All states above pass in dark mode too |

### 2.2 How to Run Visual Tests

```bash
# Run all Playwright tests
npx playwright test

# Run with UI mode for interactive debugging
npx playwright test --ui

# Run a specific test file
npx playwright test tests/components/button.spec.ts

# Update snapshots after intentional visual change
npx playwright test --update-snapshots

# Run headed (see the browser)
npx playwright test --headed

# Debug a single test
npx playwright test --debug tests/components/button.spec.ts
```

### 2.3 Playwright Test Structure

Every component test file must follow this structure:

```typescript
import { test, expect } from '@playwright/test';

test.describe('ComponentName', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/path-to-component-or-storybook-url');
    await page.waitForLoadState('networkidle');
  });

  test('renders correctly in default state', async ({ page }) => {
    await expect(page.locator('[data-testid="component"]')).toBeVisible();
    await expect(page).toHaveScreenshot('component-default.png');
  });

  test('hover state is visually correct', async ({ page }) => {
    const el = page.locator('[data-testid="component"]');
    await el.hover();
    await expect(page).toHaveScreenshot('component-hover.png');
  });

  test('focus state shows ring', async ({ page }) => {
    await page.keyboard.press('Tab');
    await expect(page).toHaveScreenshot('component-focus.png');
  });

  test('renders correctly on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await expect(page).toHaveScreenshot('component-mobile.png');
  });
});
```

### 2.4 Snapshot Policy

- Snapshots live in `tests/__snapshots__/`.
- Never commit a snapshot update without first reviewing the visual diff.
- When a snapshot intentionally changes (design update), note it explicitly in the commit message: `fix(button): update hover snapshot after design.md token change`.
- Flaky snapshots must be investigated, not just re-run. Anti-flake strategies: `waitForLoadState('networkidle')`, mask dynamic content (timestamps, avatars), disable CSS animations in test config.

### 2.5 Playwright Config Requirements

Ensure `playwright.config.ts` includes:

```typescript
use: {
  // Disable animations so hover/focus states are stable
  reducedMotion: 'reduce',
  // Consistent font rendering
  deviceScaleFactor: 1,
  // Always screenshot on failure
  screenshot: 'only-on-failure',
  // Record video on first retry
  video: 'on-first-retry',
  // Trace on first retry for debugging
  trace: 'on-first-retry',
},
```

---

## 3. Project-Wide Sweep Rule

**This is mandatory.** When fixing any visual or interactive bug, the agent must:

1. **Identify the root cause class or pattern** — e.g., "the hover background is missing because `hover:bg-*` class is absent."
2. **Search the entire frontend codebase** for the same pattern before touching anything:

```bash
# Search for the problematic pattern across all component files
grep -r "the-pattern" src/ --include="*.tsx" --include="*.ts" --include="*.jsx" --include="*.js" -l

# Example: find all elements missing hover states on interactive elements
grep -r "cursor-pointer" src/ --include="*.tsx" -l
```

3. **Fix every instance** in the same PR/commit. A partial fix that leaves broken siblings is not acceptable.
4. **Write a Playwright test** that would have caught the original bug, so it cannot regress anywhere in the project.
5. **Document the sweep** in the PR description: list every file touched and why.

### Examples of sweep-worthy problems

- An interactive element missing a `hover:` variant → sweep all interactive elements.
- A focus ring missing on a button → sweep all focusable elements.
- A transition missing on a state change → sweep all animated elements.
- A text truncation overflow bug → sweep all text containers with fixed widths.
- A shadow or border inconsistency → sweep all card/panel components.
- A z-index stacking issue → sweep all modals, dropdowns, tooltips.
- An icon misaligned vertically → sweep all icon-containing elements.

---

## 4. Interactive State Standards

Every interactive element **must** implement all of the following states. Use Playwright to verify each one exists and is visually distinct from idle.

```
idle → hover → focus-visible → active → disabled
```

- `hover:` — must provide clear visual feedback. Never rely on cursor change alone.
- `focus-visible:` — must show a visible focus ring. Do not suppress outlines without a replacement.
- `active:` — must provide a pressed/depressed visual cue.
- `disabled:` — must be visually muted (`opacity-50` or equivalent) and non-interactive (`pointer-events-none`).
- `aria-*` attributes must match the visual state. If it looks disabled, it must be `aria-disabled="true"`.

---

## 5. Layout & Overflow Integrity

Before marking any task done, verify:

- [ ] No horizontal scroll appears at any tested viewport.
- [ ] Text does not overflow its container at any tested viewport.
- [ ] Images and media have explicit aspect ratios or `object-fit` to prevent layout shift.
- [ ] Flex/grid children do not collapse to zero width/height unexpectedly.
- [ ] Fixed/sticky elements do not cover interactive content.
- [ ] Modals and overlays have proper scroll lock and do not cause body scroll bleed.
- [ ] Before writing any frontend code use your this skill "frontend-patterns".

Test overflow with Playwright:

```typescript
test('no horizontal overflow on mobile', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 812 });
  await page.goto('/');
  const scrollWidth = await page.evaluate(() => document.body.scrollWidth);
  const clientWidth = await page.evaluate(() => document.body.clientWidth);
  expect(scrollWidth).toBeLessThanOrEqual(clientWidth);
});
```

---

## 6. Accessibility Checks

Run accessibility assertions as part of every Playwright test suite:

```bash
npm install @axe-core/playwright
```

```typescript
import AxeBuilder from '@axe-core/playwright';

test('passes accessibility audit', async ({ page }) => {
  await page.goto('/target-page');
  const results = await new AxeBuilder({ page }).analyze();
  expect(results.violations).toHaveLength(0);
});
```

Minimum requirements:
- All images have meaningful `alt` text or `alt=""` if decorative.
- All form inputs have associated `<label>` elements.
- Color contrast must pass WCAG AA (do not verify by eye — use axe).
- Keyboard navigation must reach every interactive element in logical order.
- All interactive elements have an accessible name (`aria-label`, `aria-labelledby`, or visible text).

---

## 7. Performance Sanity Checks

After significant layout or asset changes, run a quick Lighthouse CI check or verify via Playwright:

```typescript
test('page renders within performance budget', async ({ page }) => {
  await page.goto('/');
  const timing = await page.evaluate(() => {
    const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
    return { lcp: nav.loadEventEnd - nav.startTime };
  });
  // Adjust threshold per page complexity
  expect(timing.lcp).toBeLessThan(3000);
});
```

- Avoid adding images without `next/image` (or equivalent) optimisation.
- Never import an entire icon library; import individual icons only.
- Lazy-load anything below the fold.

---

## 8. Component File Checklist

Before committing any component change, confirm:

- [ ] Tailwind classes only — no inline styles, no CSS modules (unless project already uses them).
- [ ] All interactive states implemented and Playwright-verified.
- [ ] `data-testid` attribute on the root element for reliable Playwright targeting.
- [ ] Snapshot tests created or updated.
- [ ] Axe accessibility test passes.
- [ ] No horizontal overflow at 375 px viewport.
- [ ] Project-wide sweep completed if the fix applies to a pattern (not just this one component).
- [ ] PR description lists all swept files.

---

## 9. Agent Workflow Summary

```
1. Read the task
2. Search the codebase to understand the full scope of the problem
3. Implement the fix using Tailwind utilities
4. Sweep the project for the same class of problem → fix all instances
5. Write / update Playwright tests covering all interactive states
6. Run: npx playwright test
7. Review any snapshot diffs — accept only if correct
8. Run axe accessibility check
9. Verify no overflow at mobile viewport
10. Commit with a descriptive message listing all files changed
```

Never skip step 6. A change that hasn't been seen running in a browser is not finished.