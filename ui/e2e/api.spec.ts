import { expect, test } from "@playwright/test";
import { capture, resetHistory } from "./helpers";

test.beforeEach(async ({ request }) => {
  await resetHistory(request);
});

test("GET /api/v1/history/{id} returns the captured record", async ({ request }) => {
  await capture(request, "/webhooks/payment", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    data: { type: "payment.created" },
  });

  const list = await request.get("/api/v1/history");
  expect(list.ok()).toBe(true);
  const records = await list.json();
  expect(records).toHaveLength(1);

  const id = records[0].id;
  expect(id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$/);

  const got = await request.get(`/api/v1/history/${id}`);
  expect(got.status()).toBe(200);

  const body = await got.json();
  expect(body.id).toBe(id);
  expect(body.request.path).toBe("/webhooks/payment");
  expect(body.request.method).toBe("POST");
});

test("GET /api/v1/history/{id} returns 404 for an unknown id", async ({ request }) => {
  const res = await request.get("/api/v1/history/does-not-exist");
  expect(res.status()).toBe(404);
});
