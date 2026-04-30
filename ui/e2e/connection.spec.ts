import { expect, test } from "@playwright/test";
import { resetHistory } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
});

test("connection indicator turns red when /healthz starts failing", async ({ page }) => {
  // First load with backend reachable — the health poll sees /healthz 200.
  await page.goto("/");

  // Wait for at least one /healthz response to establish "connected".
  await page.waitForResponse((r) => r.url().endsWith("/healthz") && r.ok(), { timeout: 5_000 });

  // Now simulate backend-down by failing all subsequent /healthz requests
  // at the browser layer. The health poll runs every 3s, so wait long enough
  // for one failed tick.
  await page.route("**/healthz", (route) => route.abort("failed"));

  // Settings page shows the status row driven by the same `listening` signal.
  await page.getByTitle("Configuration").click();

  // "Disconnected" appears in the Runtime section after the poll fails.
  await expect(page.getByText(/Disconnected/)).toBeVisible({ timeout: 8_000 });
});
