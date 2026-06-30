import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('AeolyzerSidebar', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('renders correctly in default state', async ({ page }) => {
    await expect(page.locator('text=Home').first()).toBeVisible();
    await expect(page).toHaveScreenshot('sidebar-default.png');
  });

  test('hover state is visually correct', async ({ page }) => {
    const el = page.locator('text=Agent').first();
    await el.hover();
    await expect(page).toHaveScreenshot('sidebar-hover.png');
  });

  test('renders correctly on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto('/');
    await expect(page).toHaveScreenshot('sidebar-mobile-closed.png');
  });

  test('renders correctly on mobile when open', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto('/');
    // Click the mobile hamburger menu in the mobile header
    await page.locator('.md\\:hidden > button[aria-label="Open sidebar"]').click();
    // Wait for animation
    await page.waitForTimeout(500);
    await expect(page).toHaveScreenshot('sidebar-mobile-open.png');
  });

  test('no horizontal overflow on mobile', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 812 });
    await page.goto('/');
    const scrollWidth = await page.evaluate(() => document.body.scrollWidth);
    const clientWidth = await page.evaluate(() => document.body.clientWidth);
    expect(scrollWidth).toBeLessThanOrEqual(clientWidth);
  });

  test('passes accessibility audit', async ({ page }) => {
    const results = await new AxeBuilder({ page })
      .disableRules(['region'])
      .analyze();
    expect(results.violations).toHaveLength(0);
  });
});
