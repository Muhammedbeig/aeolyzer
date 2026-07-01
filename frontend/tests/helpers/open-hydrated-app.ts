import { expect, type Page } from "@playwright/test"

export async function openHydratedApp(page: Page) {
  const hydrationRequest = page.waitForResponse((response) => {
    const url = new URL(response.url())
    return url.origin === "http://localhost:8080"
      && url.pathname === "/v1/conversations"
  })

  await page.goto("/", { waitUntil: "domcontentloaded" })
  await hydrationRequest
  await expect(page.getByText("AEOlyzer", { exact: true }).first()).toBeVisible()
}
