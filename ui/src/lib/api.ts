// Thin HTTP layer. One place to change URLs.

export const api = {
  getHistory: () => fetch("/api/v1/history"),
  clearHistory: () => fetch("/api/v1/history", { method: "DELETE" }),
  getSettings: () => fetch("/api/v1/settings"),
  getVersion: () => fetch("/version", { cache: "no-store" }),
  getHealth: () => fetch("/healthz", { cache: "no-store" }),
  subscribe: () => new EventSource("/api/v1/events"),
};
