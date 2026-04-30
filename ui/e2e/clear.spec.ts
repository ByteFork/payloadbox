import { expect, test } from "@playwright/test";
import { capture, resetHistory } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
});

test("Clear button empties the list", async ({ page, request }) => {
  await capture(request, "/x", { method: "POST" });
  await capture(request, "/y", { method: "POST" });

  await page.goto("/");
  await expect(page.locator(".req-row")).toHaveCount(2);

  await page.getByRole("button", { name: /Clear/ }).click();

  await expect(page.locator(".req-row")).toHaveCount(0);
  await expect(page.getByText(/No requests yet/)).toBeVisible();
});
