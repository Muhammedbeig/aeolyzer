import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';
import { openHydratedApp } from '../helpers/open-hydrated-app';

const profileDocument = {
  section: 'profile',
  version: 1,
  profile: {
    name: 'AEOlyzer',
    description: 'AI visibility and content intelligence platform.',
  },
  updated_at: '2026-07-01T04:00:00Z',
};

test.describe('Knowledge base', () => {
  test.beforeEach(async ({ page }) => {
    await page.clock.setFixedTime(new Date('2026-07-01T04:00:00Z'));
    await page.route('http://localhost:8080/**', async (route) => {
      const request = route.request();
      const url = new URL(request.url());
      if (
        request.method() === 'GET' &&
        url.pathname === '/v1/conversations'
      ) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ conversations: [] }),
        });
        return;
      }
      if (
        request.method() === 'GET' &&
        url.pathname === '/v1/knowledge/profile'
      ) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(profileDocument),
        });
        return;
      }
      if (
        request.method() === 'GET' &&
        url.pathname === '/v1/knowledge/memory'
      ) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            section: 'memory',
            version: 0,
            memory: { facts: [] },
          }),
        });
        return;
      }
      if (
        request.method() === 'PUT' &&
        url.pathname.startsWith('/v1/knowledge/')
      ) {
        const update = request.postDataJSON();
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            section: url.pathname.split('/').pop(),
            version: update.version + 1,
            ...update,
            approved: undefined,
            updated_at: '2026-07-01T04:01:00Z',
          }),
        });
        return;
      }
      await route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({
          error: { code: 'not_found', message: 'Not found' },
        }),
      });
    });
    await openHydratedApp(page);
  });

  test('loads and saves the profile through the backend', async ({ page }) => {
    await page.getByRole('button', { name: 'Profile' }).click();
    await expect(page.getByLabel('Company/Agent Name')).toHaveValue('AEOlyzer');
    await expect(page).toHaveScreenshot('knowledge-profile-loaded.png');

    await page
      .getByLabel('Core Description')
      .fill('AEO and SEO intelligence for modern brands.');
    const updateRequest = page.waitForRequest(
      (request) =>
        request.method() === 'PUT' &&
        request.url().endsWith('/v1/knowledge/profile'),
    );
    await page.getByRole('button', { name: 'Save Profile' }).click();
    const request = await updateRequest;
    expect(request.postDataJSON()).toMatchObject({
      version: 1,
      approved: true,
      profile: {
        name: 'AEOlyzer',
        description: 'AEO and SEO intelligence for modern brands.',
      },
    });
  });

  test('renders and updates the empty memory state', async ({ page }) => {
    await page.getByRole('button', { name: 'Memory' }).click();
    await expect(page.getByText('Memory Vault is Empty')).toBeVisible();
    await expect(page).toHaveScreenshot('knowledge-memory-empty.png');

    await page
      .getByLabel('Approved memory fact')
      .fill('Primary audience is technical founders.');
    await page.getByRole('button', { name: 'Add Fact' }).click();
    await expect(
      page.getByText('Primary audience is technical founders.'),
    ).toBeVisible();
  });

  test('keeps focus visible and has no mobile overflow', async ({ page }) => {
    await page.getByRole('button', { name: 'Profile' }).click();
    const nameInput = page.getByLabel('Company/Agent Name');
    await nameInput.focus();
    await expect(page).toHaveScreenshot('knowledge-profile-focus.png');

    await page.setViewportSize({ width: 375, height: 812 });
    const dimensions = await page.evaluate(() => ({
      scrollWidth: document.body.scrollWidth,
      clientWidth: document.body.clientWidth,
    }));
    expect(dimensions.scrollWidth).toBeLessThanOrEqual(dimensions.clientWidth);
  });

  test('passes accessibility audit', async ({ page }) => {
    await page.getByRole('button', { name: 'Profile' }).click();
    const results = await new AxeBuilder({ page })
      .disableRules(['region'])
      .analyze();
    expect(results.violations).toHaveLength(0);
  });

  test('shows a recoverable backend error', async ({ page }) => {
    await page.route(
      'http://localhost:8080/v1/knowledge/profile',
      async (route) => {
        await route.fulfill({
          status: 503,
          contentType: 'application/json',
          body: JSON.stringify({
            error: {
              code: 'temporarily_unavailable',
              message: 'Knowledge settings are temporarily unavailable.',
            },
          }),
        });
      },
    );

    await page.getByRole('button', { name: 'Profile' }).click();
    await expect(page.getByTestId('knowledge-error')).toContainText(
      'Knowledge settings are temporarily unavailable.',
    );
    await expect(page.getByRole('button', { name: 'Retry' })).toBeVisible();
  });
});
