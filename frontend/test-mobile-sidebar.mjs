import { chromium } from 'playwright';

(async () => {
  console.log('Starting Playwright...');
  const browser = await chromium.launch();
  const context = await browser.newContext({
    viewport: { width: 375, height: 667 }, // iPhone SE dimensions
    deviceScaleFactor: 2,
    isMobile: true,
    hasTouch: true
  });

  const page = await context.newPage();
  
  console.log('Navigating to localhost:3000...');
  try {
    await page.goto('http://localhost:3000', { waitUntil: 'networkidle', timeout: 15000 });
  } catch (err) {
    console.error('Failed to load page. Is the dev server running?', err);
    await browser.close();
    process.exit(1);
  }

  console.log('Page loaded. Taking initial screenshot (Sidebar should be closed)...');
  await page.screenshot({ path: 'sidebar-mobile-initial.png' });

  console.log('Clicking the menu button to open sidebar...');
  // The menu button is an SVG wrapped in a button inside the mobile header.
  // We can look for the button containing the Menu icon, or simply the first button in the mobile header
  await page.locator('button[aria-label="Open sidebar"]').click();
  
  // Wait for sidebar transition (300ms + some buffer)
  await page.waitForTimeout(500);

  console.log('Taking screenshot of open sidebar...');
  await page.screenshot({ path: 'sidebar-mobile-opened.png' });

  console.log('Clicking the overlay to close the sidebar...');
  // The overlay is the fixed div with bg-black/50 that appears when mobile sidebar is open
  await page.locator('.bg-black\\/50').click();
  
  // Wait for sidebar transition
  await page.waitForTimeout(500);

  console.log('Taking screenshot of closed sidebar...');
  await page.screenshot({ path: 'sidebar-mobile-closed.png' });

  console.log('Done!');
  await browser.close();
})();
