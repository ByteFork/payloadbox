import type { APIRequestContext, Page } from "@playwright/test";

// Clear captured history between tests. Backend is single-instance with
// in-memory state, so tests share it unless reset.
export async function resetHistory(request: APIRequestContext): Promise<void> {
  await request.delete("/api/v1/history");
}

// Send a captured request (what webhook senders would do in real use).
export async function capture(
  request: APIRequestContext,
  path: string,
  opts?: { method?: string; headers?: Record<string, string>; data?: unknown },
): Promise<void> {
  const method = opts?.method ?? "POST";
  await request.fetch(path, {
    method,
    headers: opts?.headers,
    data: opts?.data,
  });
}

// Wait for at least one row to appear in the request list.
export async function waitForRow(page: Page): Promise<void> {
  await page.locator(".req-row").first().waitFor({ state: "visible" });
}
