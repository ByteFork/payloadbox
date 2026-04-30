import { expect, test } from "@playwright/test";
import { capture, resetHistory, waitForRow } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
});

test("POST is captured, listed, and details render in the inspector", async ({ page, request }) => {
  await capture(request, "/webhooks/payment", {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Test": "hello" },
    data: { type: "payment.created", amount: 4200 },
  });

  await page.goto("/");
  await waitForRow(page);

  // Row shows the method and path.
  const row = page.locator(".req-row").first();
  await expect(row).toContainText("POST");
  await expect(row).toContainText("/webhooks/payment");

  // Overview is the default tab and shows the stat cards.
  await row.click();
  await expect(page.getByText("Status", { exact: true })).toBeVisible();
  await expect(page.getByText("Duration", { exact: true })).toBeVisible();
  await expect(page.getByText("Response Body", { exact: true })).toBeVisible();

  // Headers tab lists the custom header we sent.
  await page.getByRole("tab", { name: /Headers/ }).click();
  await expect(page.getByText("X-Test")).toBeVisible();
  await expect(page.getByText("hello")).toBeVisible();

  // Body tab shows the payload JSON.
  await page.getByRole("tab", { name: /^Body/ }).click();
  await expect(page.locator("pre").first()).toContainText("payment.created");

  // cURL tab renders a curl command that faithfully preserves method, path,
  // headers, and body — the whole point of the replay feature.
  await page.getByRole("tab", { name: /cURL/ }).click();
  const curl = page.locator("pre").last();
  await expect(curl).toContainText("curl -X POST");
  await expect(curl).toContainText("/webhooks/payment");
  await expect(curl).toContainText("-H 'Content-Type: application/json'");
  await expect(curl).toContainText("-H 'X-Test: hello'");
  await expect(curl).toContainText("payment.created");
});

test("body that exceeds MAX_BODY_SIZE_BYTES records a 413 attempt", async ({ page, request }) => {
  // Default cap is 5120 bytes; send well over.
  const huge = "x".repeat(8192);
  await capture(request, "/big", {
    headers: { "Content-Type": "text/plain" },
    data: huge,
  });

  await page.goto("/");
  await waitForRow(page);

  const row = page.locator(".req-row").first();
  await expect(row).toContainText("413");
});
