export function formatDuration(ns: number | undefined | null): string {
  if (ns === undefined || ns === null) return "—";
  const ms = ns / 1e6;
  if (ms < 1) return `${(ns / 1e3).toFixed(0)}µs`;
  if (ms < 1000) return `${ms.toFixed(ms < 10 ? 1 : 0)}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
}

export function formatBytes(n: number | undefined | null): string {
  if (!n) return "-";
  if (n < 1024) return `${n}b`;
  if (n < 1024 * 1024) return `${parseFloat((n / 1024).toFixed(1))}kb`;
  if (n < 1024 * 1024 * 1024) return `${parseFloat((n / (1024 * 1024)).toFixed(2))}mb`;
  return `${parseFloat((n / (1024 * 1024 * 1024)).toFixed(2))}gb`;
}

export function formatTime(ts: string | number | undefined | null): string {
  if (!ts) return "—";
  const d = new Date(ts);
  return d.toLocaleTimeString("en-US", {
    hour12: false,
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

export function formatRelTime(ts: string | number | undefined | null): string {
  if (!ts) return "—";
  const diff = Date.now() - new Date(ts).getTime();
  if (diff < 5000) return "just now";
  if (diff < 60000) return `${Math.floor(diff / 1000)}s ago`;
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`;
  return `${Math.floor(diff / 3600000)}h ago`;
}

export function formatFullTime(ts: string | number | undefined | null): string {
  if (!ts) return "—";
  return new Date(ts).toLocaleString();
}

const METHOD_COLORS: Record<string, { color: string; bg: string }> = {
  GET: { color: "#2563eb", bg: "var(--get-bg)" },
  POST: { color: "#059669", bg: "var(--post-bg)" },
  PUT: { color: "#d97706", bg: "var(--put-bg)" },
  PATCH: { color: "#7c3aed", bg: "var(--patch-bg)" },
  DELETE: { color: "#dc2626", bg: "var(--delete-bg)" },
  HEAD: { color: "#0891b2", bg: "var(--head-bg)" },
  OPTIONS: { color: "#65a30d", bg: "var(--options-bg)" },
};

export function methodInfo(method: string): { color: string; bg: string } {
  return METHOD_COLORS[method.toUpperCase()] || { color: "#6b7280", bg: "var(--surface2)" };
}

export function statusInfo(code: number | undefined | null): { color: string; bg: string } {
  if (!code) return { color: "var(--text-4)", bg: "var(--surface2)" };
  if (code < 300) return { color: "var(--ok)", bg: "var(--ok-bg)" };
  if (code < 400) return { color: "var(--info)", bg: "var(--info-bg)" };
  if (code < 500) return { color: "var(--warn)", bg: "var(--warn-bg)" };
  return { color: "var(--err)", bg: "var(--err-bg)" };
}

export function firstHeader(headers: Record<string, string[]> | undefined, name: string): string {
  if (!headers) return "";
  const lower = name.toLowerCase();
  for (const k of Object.keys(headers)) {
    if (k.toLowerCase() === lower) {
      const v = headers[k];
      return Array.isArray(v) ? (v[0] ?? "") : String(v);
    }
  }
  return "";
}

export function flattenHeaders(headers: Record<string, string[]> | undefined): [string, string][] {
  if (!headers) return [];
  return Object.entries(headers).map(([k, v]) => [k, Array.isArray(v) ? v.join(", ") : String(v)]);
}

export function highlightJSON(src: string | undefined | null): string {
  if (!src) return "";
  let pretty = src;
  try {
    pretty = JSON.stringify(JSON.parse(src), null, 2);
  } catch {
    /* leave as-is */
  }
  return pretty.replace(
    /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+-]?\d+)?)/g,
    (match) => {
      let cls = "json-num";
      if (match.startsWith('"')) cls = match.endsWith(":") ? "json-key" : "json-str";
      else if (/true|false/.test(match)) cls = "json-bool";
      else if (/null/.test(match)) cls = "json-null";
      return `<span class="${cls}">${match}</span>`;
    },
  );
}

// ── shell / yaml code highlighters ───────────────────────────────────────────
// Colors are for dark code-block backgrounds (var(--code-bg)).
const SHELL_C = {
  comment: "#4b5563",
  cmd: "#a78bfa",
  flag: "#94a3b8",
  method: "#fbbf24",
  str: "#fbbf24",
  url: "#60a5fa",
  lineCont: "#64748b",
};

function escHtml(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

export function highlightShell(code: string): string {
  return code
    .split("\n")
    .map((line) => {
      if (/^\s*#/.test(line))
        return `<span style="color:${SHELL_C.comment}">${escHtml(line)}</span>`;
      return escHtml(line)
        .replace(
          /\b(curl|chmod|open|docker|tar|\.\/payloadbox)\b/g,
          `<span style="color:${SHELL_C.cmd}">$1</span>`,
        )
        .replace(
          /\b(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\b/g,
          `<span style="color:${SHELL_C.method}">$1</span>`,
        )
        .replace(/(?<![:\w])(-[A-Za-z]+)/g, `<span style="color:${SHELL_C.flag}">$1</span>`)
        .replace(/'([^']*)'/g, `<span style="color:${SHELL_C.str}">'$1'</span>`)
        .replace(/(https?:\/\/[^\s<]+)/g, `<span style="color:${SHELL_C.url}">$1</span>`)
        .replace(/\\$/, `<span style="color:${SHELL_C.lineCont}">\\</span>`);
    })
    .join("\n");
}
