import { defineConfig } from "@playwright/test";

// Expects the Go binary at repo root: run `cd ui && pnpm build && cd .. && go build`
// before `pnpm test:e2e` locally.
export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false, // backend has shared in-memory state
  workers: 1,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? "github" : "list",

  use: {
    baseURL: "http://localhost:8080",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },

  webServer: {
    command: "../payloadbox",
    url: "http://localhost:8080/healthz",
    reuseExistingServer: !process.env.CI,
    timeout: 15_000,
    stdout: "pipe",
    stderr: "pipe",
  },

  projects: [{ name: "chromium", use: { browserName: "chromium" } }],
});
