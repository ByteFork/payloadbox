import { expect, test } from "@playwright/test";
import { capture, resetHistory } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
  await capture(request, "/api/users", { method: "GET" });
  await capture(request, "/api/orders", { method: "POST" });
  await capture(request, "/webhooks/hook", { method: "POST" });
  await capture(request, "/api/session", { method: "DELETE" });
});

test("search narrows the list and empty-state appears when nothing matches", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator(".req-row")).toHaveCount(4);

  await page.getByPlaceholder("Filter requests…").fill("webhooks");
  await expect(page.locator(".req-row")).toHaveCount(1);
  await expect(page.locator(".req-row").first()).toContainText("/webhooks/hook");

  await page.getByPlaceholder("Filter requests…").fill("nomatch-xyz");
  await expect(page.locator(".req-row")).toHaveCount(0);
  await expect(page.getByText(/No requests match your filter/)).toBeVisible();

  // Clearing the input brings all rows back.
  await page.getByPlaceholder("Filter requests…").fill("");
  await expect(page.locator(".req-row")).toHaveCount(4);
});

test("method chips filter by HTTP method", async ({ page }) => {
  await page.goto("/");
  await expect(page.locator(".req-row")).toHaveCount(4);

  // Scope to the filter bar via the search input's parent; the filter row
  // holds buttons whose textContent is exactly the method name. MethodBadge
  // spans inside rows use the same text but aren't buttons, so getByRole
  // won't collide.
  await page.getByRole("button", { name: "POST", exact: true }).click();
  await expect(page.locator(".req-row")).toHaveCount(2);
  for (const row of await page.locator(".req-row").all()) {
    await expect(row).toContainText("POST");
  }

  await page.getByRole("button", { name: "DELETE", exact: true }).click();
  await expect(page.locator(".req-row")).toHaveCount(1);
  await expect(page.locator(".req-row").first()).toContainText("/api/session");

  await page.getByRole("button", { name: "ALL", exact: true }).click();
  await expect(page.locator(".req-row")).toHaveCount(4);
});

test("search and method chip compose", async ({ page }) => {
  await page.goto("/");

  await page.getByRole("button", { name: "POST", exact: true }).click();
  await expect(page.locator(".req-row")).toHaveCount(2);

  await page.getByPlaceholder("Filter requests…").fill("webhooks");
  await expect(page.locator(".req-row")).toHaveCount(1);
  await expect(page.locator(".req-row").first()).toContainText("/webhooks/hook");
});
