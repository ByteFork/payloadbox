import { expect, test } from "@playwright/test";
import { capture, resetHistory, waitForRow } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
});

test("new captures appear without reload (SSE)", async ({ page, request }) => {
  await page.goto("/");

  // Empty state first.
  await expect(page.getByText(/No requests yet/)).toBeVisible();

  // Send the capture with the page already open. No navigation after this.
  await capture(request, "/live", { method: "POST", data: { ping: 1 } });

  await waitForRow(page);
  await expect(page.locator(".req-row").first()).toContainText("/live");
});

test("multiple SSE events accumulate in order (newest first)", async ({ page, request }) => {
  await page.goto("/");

  for (const p of ["/a", "/b", "/c"]) {
    await capture(request, p, { method: "POST" });
  }

  // Wait until three rows appear.
  await expect(page.locator(".req-row")).toHaveCount(3);

  // Newest first in the DOM.
  const rows = page.locator(".req-row");
  await expect(rows.nth(0)).toContainText("/c");
  await expect(rows.nth(1)).toContainText("/b");
  await expect(rows.nth(2)).toContainText("/a");
});
