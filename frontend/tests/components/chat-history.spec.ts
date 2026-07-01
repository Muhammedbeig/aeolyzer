import { test, expect } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { openHydratedApp } from '../helpers/open-hydrated-app';

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
  content_type: 'product_description',
  title: 'Product launch article',
  starred: true,
  created_at: '2026-06-30T11:00:00Z',
  updated_at: '2026-06-30T11:00:00Z',
};

test.describe('Chat history and attachments', () => {
  test.beforeEach(async ({ page }) => {
    await page.clock.setFixedTime(new Date('2026-07-01T04:00:00Z'));
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
      if (request.method() === 'POST' && url.pathname === '/v1/conversations') {
        const input = request.postDataJSON();
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'new-content-conversation',
            agent: input.agent,
            content_type: input.content_type,
            title: 'New chat',
            starred: false,
            created_at: '2026-07-01T04:00:00Z',
            updated_at: '2026-07-01T04:00:00Z',
          }),
        });
        return;
      }
      if (
        request.method() === 'POST' &&
        url.pathname.endsWith('/messages')
      ) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            conversation: {
              id: 'new-content-conversation',
              agent: 'content',
              content_type: 'blog_post',
              title: 'Draft a launch story',
              starred: false,
              created_at: '2026-07-01T04:00:00Z',
              updated_at: '2026-07-01T04:01:00Z',
            },
            user_message: {
              id: 'new-user-message',
              role: 'user',
              content: 'Draft a launch story',
              created_at: '2026-07-01T04:01:00Z',
            },
            reply: {
              id: 'new-assistant-message',
              role: 'assistant',
              content: 'I will draft this as a blog post.',
              created_at: '2026-07-01T04:01:01Z',
            },
          }),
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
    await openHydratedApp(page);
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

  test('shows removable file previews without composer borders', async ({ page }) => {
    await page.getByText('Agent', { exact: true }).click();
    const fileInput = page.getByLabel('Choose attachments');
    await fileInput.setInputFiles([
      {
        name: 'evidence.txt',
        mimeType: 'text/plain',
        buffer: Buffer.from('safe evidence'),
      },
      {
        name: 'audit.pdf',
        mimeType: 'application/pdf',
        buffer: Buffer.from('%PDF-1.7\n%%EOF'),
      },
      {
        name: 'diagram.png',
        mimeType: 'image/png',
        buffer: readFileSync(
          'tests/components/sidebar.spec.ts-snapshots/sidebar-mobile-closed-chromium-win32.png',
        ),
      },
    ]);
    const sendButton = page.getByRole('button', { name: 'Send message' });
    await expect(sendButton).toBeEnabled();
    await expect(page.getByTestId('attachment-preview-card')).toHaveCount(3);
    const imagePreview = page.getByAltText('Preview of diagram.png');
    await expect(imagePreview).toBeVisible();
    await expect
      .poll(() => imagePreview.evaluate((image) => (image as HTMLImageElement).naturalWidth))
      .toBeGreaterThan(0);
    await expect(
      page.getByRole('button', { name: /Add attachment, 3 selected/ }),
    ).toBeVisible();
    const composer = page.getByTestId('chat-input').locator('form > div');
    await expect(composer).toHaveCSS('border-top-width', '0px');
    await expect(page).toHaveScreenshot('attachment-previews.png');

    await page.getByRole('button', { name: 'Remove evidence.txt' }).click();
    await expect(page.getByTestId('attachment-preview-card')).toHaveCount(2);

    await page.setViewportSize({ width: 375, height: 812 });
    await expect(page).toHaveScreenshot('attachment-previews-mobile.png');
  });

  test('sends the selected content type to the backend', async ({ page }) => {
    await page.getByText('Content', { exact: true }).click();
    await page.getByRole('button', { name: 'Blog Post' }).click();
    await expect(
      page.getByRole('button', { name: 'Blog Post' }),
    ).toHaveAttribute('aria-pressed', 'true');
    await page
      .getByPlaceholder('Describe what you want to write...')
      .fill('Draft a launch story');

    const messageRequest = page.waitForRequest(
      (request) =>
        request.method() === 'POST' &&
        request.url().endsWith('/new-content-conversation/messages'),
    );
    await page.getByRole('button', { name: 'Send message' }).click();
    const request = await messageRequest;
    expect(request.postData()).toContain('blog_post');
    await expect(page.getByText('I will draft this as a blog post.')).toBeVisible();
  });
});
