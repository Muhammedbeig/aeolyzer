import { test, expect } from '@playwright/test';

const auditConversation = {
  id: 'audit-conversation',
  agent: 'audit',
  title: 'Homepage visibility audit',
  starred: false,
  created_at: '2026-06-30T10:00:00Z',
  updated_at: '2026-06-30T10:00:00Z',
};

const contentConversation = {
  id: 'content-conversation',
  agent: 'content',
  title: 'Product launch article',
  starred: true,
  created_at: '2026-06-30T11:00:00Z',
  updated_at: '2026-06-30T11:00:00Z',
};

test.describe('Chat history and attachments', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('http://localhost:8080/**', async (route) => {
      const request = route.request();
      const url = new URL(request.url());
      if (request.method() === 'GET' && url.pathname === '/v1/conversations') {
        const conversations =
          url.searchParams.get('agent') === 'content'
            ? [contentConversation]
            : [auditConversation];
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ conversations }),
        });
        return;
      }
      if (
        request.method() === 'GET' &&
        url.pathname.endsWith('/messages')
      ) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            messages: [
              {
                id: 'user-message',
                role: 'user',
                content: 'Review the attached evidence',
                created_at: '2026-06-30T10:01:00Z',
              },
              {
                id: 'assistant-message',
                role: 'assistant',
                content: 'The evidence is ready for review.',
                created_at: '2026-06-30T10:01:01Z',
              },
            ],
          }),
        });
        return;
      }
      await route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({ error: { code: 'not_found', message: 'Not found' } }),
      });
    });
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('keeps Audit and Content histories isolated', async ({ page }) => {
    await page.getByText('Agent', { exact: true }).click();
    await expect(page.getByText(auditConversation.title)).toBeVisible();
    await expect(page.getByText(contentConversation.title)).toHaveCount(0);

    await page.getByText('Content', { exact: true }).click();
    await expect(page.getByText(contentConversation.title)).toBeVisible();
    await expect(page.getByText(auditConversation.title)).toHaveCount(0);
    await expect(page).toHaveScreenshot('content-history-isolated.png');
  });

  test('resumes a stored conversation', async ({ page }) => {
    await page.getByText('Agent', { exact: true }).click();
    await page.getByText(auditConversation.title).click();
    await expect(page.getByText('Review the attached evidence')).toBeVisible();
    await expect(page.getByText('The evidence is ready for review.')).toBeVisible();
    await expect(page).toHaveScreenshot('audit-conversation-resumed.png');
  });

  test('attachment-only input enables send without changing the composer layout', async ({ page }) => {
    await page.getByText('Agent', { exact: true }).click();
    const fileInput = page.getByLabel('Choose attachments');
    await fileInput.setInputFiles({
      name: 'evidence.txt',
      mimeType: 'text/plain',
      buffer: Buffer.from('safe evidence'),
    });
    const sendButton = page.getByRole('button', { name: 'Send message' });
    await expect(sendButton).toBeEnabled();
    await expect(page.getByRole('button', { name: /Add attachment, 1 selected/ })).toBeVisible();
  });
});
