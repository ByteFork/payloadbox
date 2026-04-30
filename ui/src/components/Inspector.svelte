<script lang="ts">
import {
  firstHeader,
  flattenHeaders,
  formatBytes,
  formatDuration,
  formatFullTime,
  formatTime,
  highlightJSON,
  highlightShell,
  methodInfo,
  statusInfo,
} from "../lib/format";
import type { RequestRecord } from "../types";
import CopyBtn from "./CopyBtn.svelte";
import MethodBadge from "./MethodBadge.svelte";
import StatusBadge from "./StatusBadge.svelte";

export let record: RequestRecord | null = null;

type Tab = "overview" | "body" | "headers" | "query" | "response" | "curl";
let tab: Tab = "overview";
let prevId = "";

$: if (record && record.id !== prevId) {
  tab = "overview";
  prevId = record.id;
}

$: req = record?.request;
$: res = record?.response;
$: mc = methodInfo(req?.method ?? "");
$: sc = res?.status_code;
$: si = statusInfo(sc);
$: totalMs = (record?.duration_ns ?? 0) / 1e6;
$: reqHeaders = flattenHeaders(req?.headers);
$: resHeaders = flattenHeaders(res?.headers);
$: reqContentType = firstHeader(req?.headers, "Content-Type") || "request body";

function buildCurl(r: RequestRecord | null): string {
  if (!r) return "";
  const q = r.request?.query ? `?${r.request.query}` : "";
  const host = firstHeader(r.request?.headers, "Host") || r.request?.host || "localhost:8080";
  const proto = firstHeader(r.request?.headers, "X-Forwarded-Proto") || "http";
  let cmd = `curl -X ${r.request.method.toUpperCase()} '${proto}://${host}${r.request.path}${q}'`;
  for (const [k, v] of flattenHeaders(r.request?.headers)) {
    cmd += ` \\\n  -H '${k}: ${v}'`;
  }
  if (r.request?.body) cmd += ` \\\n  -d '${r.request.body.replace(/'/g, "'\\''")}'`;
  return cmd;
}

$: tabs = [
  { id: "overview" as Tab, label: "Overview", badge: null },
  { id: "body" as Tab, label: "Body", badge: req?.body ? "●" : null },
  { id: "headers" as Tab, label: "Headers", badge: reqHeaders.length },
  { id: "query" as Tab, label: "Query", badge: req?.query ? req.query.split("&").length : null },
  { id: "response" as Tab, label: "Response", badge: sc ?? null },
  { id: "curl" as Tab, label: "cURL", badge: null },
];

$: sizeFormatted = formatBytes(res?.size_in_bytes);
$: responseColor = sizeFormatted === "—" ? "var(--text-4)" : "var(--info)";
$: durationColor = totalMs > 200 ? "var(--err)" : totalMs > 50 ? "var(--warn)" : "var(--ok)";
$: durationLabel = totalMs > 200 ? "Slow" : totalMs > 50 ? "Moderate" : "Fast";
</script>

{#if !record || !req}
  <div
    class="flex-1 flex flex-col items-center justify-center gap-4"
    style="background: var(--bg);"
  >
    <div
      class="flex items-center justify-center rounded-full"
      style="width: 56px; height: 56px; background: var(--surface); border: 1px solid var(--border);"
    >
      <span class="material-symbols-outlined" style="font-size: 24px; color: var(--text-4);">arrow_back</span>
    </div>
    <div class="text-center">
      <div style="font-size: var(--fs-body); font-weight: 500; color: var(--text-3); margin-bottom: 4px;">Select a request</div>
      <div style="font-size: var(--fs-ui); color: var(--text-4);">Click any request to inspect its details</div>
    </div>
  </div>
{:else}
  <div
    class="flex-1 flex flex-col detail-enter"
    style="background: var(--bg); min-width: 0; overflow: hidden;"
  >
    <!-- Header -->
    <div style="background: var(--surface); border-bottom: 1px solid var(--border); padding: 14px 24px;">
      <div class="flex items-center gap-2.5 flex-wrap gap-y-2" style="margin-bottom: 8px;">
        <MethodBadge method={req.method} />
        {#if sc}<StatusBadge code={sc} text={res?.status_text} />{/if}
        <div class="flex items-center gap-1.5 ml-auto flex-wrap gap-y-2">
          <span class="font-mono" style="font-size: var(--fs-data); color: var(--text-4);">{formatTime(record.created_at)}</span>
          <div style="width: 1px; height: 12px; background: var(--border-strong);"></div>
          <span class="font-mono" style="font-size: var(--fs-data); color: var(--text-4);">{formatDuration(record.duration_ns)}</span>
          <div style="width: 1px; height: 12px; background: var(--border-strong);"></div>
          <CopyBtn text={buildCurl(record)} label="cURL" />
        </div>
      </div>
      <div
        class="font-mono truncate"
        style="font-size: var(--fs-emph); color: var(--text);"
      >
        <span style="color: {mc.color}; font-weight: 600;">{req.path}</span>{#if req.query}<span style="color: var(--text-3);">?{req.query}</span>{/if}
      </div>
      <div class="flex items-center gap-2" style="margin-top: 6px;">
        <span class="material-symbols-outlined" style="font-size: 12px; color: var(--text-4);">wifi</span>
        <span class="font-mono" style="font-size: var(--fs-data); color: var(--text-4);">{req.remote_addr ?? "—"}</span>
        <span style="color: var(--text-4);">·</span>
        <span class="material-symbols-outlined" style="font-size: 12px; color: var(--text-4);">data_usage</span>
        <span class="font-mono" style="font-size: var(--fs-data); color: var(--text-4);">{sizeFormatted}</span>
      </div>
    </div>

    <!-- Tabs -->
    <div
      role="tablist"
      class="flex gap-0.5"
      style="background: var(--surface); border-bottom: 1px solid var(--border); padding: 0 16px;
             overflow-x: auto; overflow-y: hidden;"
    >
      {#each tabs as t}
        {@const active = tab === t.id}
        <button
          type="button"
          role="tab"
          aria-selected={active}
          class="flex items-center gap-1.5 whitespace-nowrap transition-all"
          on:click={() => (tab = t.id)}
          style="padding: 10px 12px; font-size: var(--fs-ui);
                 font-weight: {active ? 600 : 400};
                 color: {active ? 'var(--primary)' : 'var(--text-3)'};
                 border-bottom: 2px solid {active ? 'var(--primary)' : 'transparent'};"
        >
          {t.label}
          {#if t.badge !== null && t.badge !== undefined}
            <span
              class="font-mono font-bold leading-[1.6]"
              style="font-size: 10px;
                     background: {active ? 'var(--primary-dim)' : 'var(--surface2)'};
                     color: {active ? 'var(--primary)' : 'var(--text-4)'};
                     border-radius: 8px; padding: 1px 5px;"
            >{t.badge}</span>
          {/if}
        </button>
      {/each}
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto" style="padding: 24px;">
      {#if tab === "overview"}
        <div class="flex flex-col gap-5">
          <!-- Stat cards -->
          <div class="grid gap-3" style="grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));">
            {#each [
              { label: "Status",        value: sc ? `${sc}` : "—", sub: res?.status_text || "—", color: si.color },
              { label: "Duration",      value: formatDuration(record.duration_ns), sub: durationLabel, color: durationColor },
              { label: "Response Body", value: sizeFormatted, sub: "size in bytes", color: responseColor },
              { label: "Headers",       value: String(reqHeaders.length), sub: "sent", color: "var(--text-2)" },
            ] as stat}
              <div
                style="background: var(--surface); border: 1px solid var(--border);
                       border-radius: 10px; padding: 14px 16px;"
              >
                <div
                  class="uppercase"
                  style="font-size: var(--fs-micro); color: var(--text-4); letter-spacing: 0.07em; margin-bottom: 6px;"
                >{stat.label}</div>
                <div
                  class="font-mono font-bold"
                  style="font-size: var(--fs-stat); color: {stat.color}; line-height: 1;"
                >{stat.value}</div>
                <div style="font-size: var(--fs-data); color: var(--text-4); margin-top: 4px;">{stat.sub}</div>
              </div>
            {/each}
          </div>

          <!-- Request Info -->
          <div
            style="background: var(--surface); border: 1px solid var(--border);
                   border-radius: 10px; padding: 20px;"
          >
            <div
              class="uppercase font-semibold"
              style="font-size: var(--fs-ui); color: var(--text-3); letter-spacing: 0.05em; margin-bottom: 14px;"
            >Request Info</div>
            <div class="flex flex-col">
              {#each [
                ["Method", req.method],
                ["Path", req.path],
                ["Query String", req.query || "—"],
                ["Remote Address", req.remote_addr || "—"],
                ["Received At", formatFullTime(record.created_at)],
                ["Content-Type", firstHeader(req.headers, "Content-Type") || "—"],
                ["User-Agent", firstHeader(req.headers, "User-Agent") || "—"],
              ] as [k, v]}
                <div
                  class="flex gap-4"
                  style="padding: 8px 0; border-bottom: 1px solid var(--border);"
                >
                  <span
                    class="shrink-0"
                    style="font-size: var(--fs-data); color: var(--text-4); width: 120px; font-weight: 500;"
                  >{k}</span>
                  <span
                    class="font-mono break-all"
                    style="font-size: var(--fs-data); color: var(--text-2);"
                  >{v}</span>
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/if}

      {#if tab === "body"}
        {#if req.body && req.body.trim() !== ""}
          <div
            style="background: var(--surface); border: 1px solid var(--border);
                   border-radius: 10px; overflow: hidden;"
          >
            <div
              class="flex items-center justify-between"
              style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
            >
              <span
                class="font-mono font-semibold"
                style="font-size: var(--fs-data); color: var(--text-3);"
              >{reqContentType.toLowerCase()}</span>
              <CopyBtn text={req.body} />
            </div>
            <div style="padding: 16px;">
              <pre
                class="font-mono"
                style="font-size: var(--fs-ui); line-height: 1.7; color: var(--text-2);"
              >{@html highlightJSON(req.body)}</pre>
            </div>
          </div>
        {:else}
          <div class="flex flex-col items-center gap-2.5" style="padding: 48px;">
            <span class="material-symbols-outlined" style="font-size: 36px; color: var(--text-4); opacity: 0.4;">indeterminate_check_box</span>
            <span style="font-size: var(--fs-ui); color: var(--text-4);">No request body</span>
          </div>
        {/if}
      {/if}

      {#if tab === "headers"}
        <div
          style="background: var(--surface); border: 1px solid var(--border);
                 border-radius: 10px; overflow: hidden;"
        >
          <div
            class="flex items-center justify-between"
            style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
          >
            <span
              class="uppercase font-semibold"
              style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
            >{reqHeaders.length} Headers</span>
            <CopyBtn text={reqHeaders.map(([k, v]) => `${k}: ${v}`).join("\n")} />
          </div>
          {#each reqHeaders as [k, v], i}
            <div
              class="flex"
              style={i < reqHeaders.length - 1 ? "border-bottom: 1px solid var(--border);" : ""}
            >
              <div
                class="shrink-0"
                style="width: 220px; padding: 10px 16px; background: var(--surface2); border-right: 1px solid var(--border);"
              >
                <span
                  class="font-mono"
                  style="font-size: var(--fs-data); color: var(--text-3); font-weight: 500;"
                >{k}</span>
              </div>
              <div class="flex-1" style="padding: 10px 16px;">
                <span
                  class="font-mono break-all"
                  style="font-size: var(--fs-data); color: var(--text-2);"
                >{v}</span>
              </div>
            </div>
          {/each}
        </div>
      {/if}

      {#if tab === "query"}
        {#if req.query}
          <div
            style="background: var(--surface); border: 1px solid var(--border);
                   border-radius: 10px; overflow: hidden;"
          >
            <div
              style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
            >
              <span
                class="uppercase font-semibold"
                style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
              >Query Parameters</span>
            </div>
            {#each req.query.split("&") as p}
              {@const [k, v] = p.split("=")}
              <div class="flex" style="border-bottom: 1px solid var(--border);">
                <div
                  class="shrink-0"
                  style="width: 200px; padding: 10px 16px; background: var(--surface2); border-right: 1px solid var(--border);"
                >
                  <span
                    class="font-mono"
                    style="font-size: var(--fs-data); color: var(--text-3); font-weight: 500;"
                  >{decodeURIComponent(k || "")}</span>
                </div>
                <div style="padding: 10px 16px;">
                  <span
                    class="font-mono break-all"
                    style="font-size: var(--fs-data); color: var(--ok);"
                  >{v ? decodeURIComponent(v) : "—"}</span>
                </div>
              </div>
            {/each}
            <div
              style="padding: 8px 16px; background: var(--surface2); border-top: 1px solid var(--border);"
            >
              <span
                class="font-mono"
                style="font-size: var(--fs-micro); color: var(--text-4);"
              >Raw: ?{req.query}</span>
            </div>
          </div>
        {:else}
          <div class="flex flex-col items-center gap-2.5" style="padding: 48px;">
            <span class="material-symbols-outlined" style="font-size: 36px; color: var(--text-4); opacity: 0.4;">manage_search</span>
            <span style="font-size: var(--fs-ui); color: var(--text-4);">No query parameters</span>
          </div>
        {/if}
      {/if}

      {#if tab === "response"}
        <div class="flex flex-col gap-4">
          {#if res}
            <div class="grid gap-3" style="grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));">
              {#each [
                { label: "Status Code", value: res.status_code ? String(res.status_code) : "—", color: si.color },
                { label: "Status Text", value: res.status_text || "—", color: "var(--text-2)" },
                { label: "Body Size", value: sizeFormatted, color: responseColor },
              ] as s}
                <div
                  style="background: var(--surface); border: 1px solid var(--border);
                         border-radius: 10px; padding: 14px 16px;"
                >
                  <div
                    class="uppercase"
                    style="font-size: var(--fs-micro); color: var(--text-4); letter-spacing: 0.07em; margin-bottom: 6px;"
                  >{s.label}</div>
                  <div
                    class="font-mono font-bold"
                    style="font-size: var(--fs-stat-sm); color: {s.color};"
                  >{s.value}</div>
                </div>
              {/each}
            </div>

            {#if resHeaders.length > 0}
              <div
                style="background: var(--surface); border: 1px solid var(--border);
                       border-radius: 10px; overflow: hidden;"
              >
                <div
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response Headers</span>
                </div>
                {#each resHeaders as [k, v]}
                  <div class="flex" style="border-bottom: 1px solid var(--border);">
                    <div
                      class="shrink-0"
                      style="width: 200px; padding: 8px 16px; background: var(--surface2); border-right: 1px solid var(--border);"
                    >
                      <span
                        class="font-mono"
                        style="font-size: var(--fs-data); color: var(--text-3); font-weight: 500;"
                      >{k}</span>
                    </div>
                    <div style="padding: 8px 16px;">
                      <span
                        class="font-mono break-all"
                        style="font-size: var(--fs-data); color: var(--text-2);"
                      >{v}</span>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}

            {#if res.body && res.body.trim() !== ""}
              <div
                style="background: var(--surface); border: 1px solid var(--border);
                       border-radius: 10px; overflow: hidden;"
              >
                <div
                  class="flex items-center justify-between"
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response Body</span>
                  <CopyBtn text={res.body} />
                </div>
                <div style="padding: 16px;">
                  <pre
                    class="font-mono"
                    style="font-size: var(--fs-ui); line-height: 1.7; color: var(--text-2);"
                  >{@html highlightJSON(res.body)}</pre>
                </div>
              </div>
            {/if}
          {:else}
            <div class="flex flex-col items-center gap-2.5" style="padding: 48px;">
              <span class="material-symbols-outlined" style="font-size: 36px; color: var(--text-4); opacity: 0.4;">reply</span>
              <span style="font-size: var(--fs-ui); color: var(--text-4);">No response data</span>
            </div>
          {/if}
        </div>
      {/if}

      {#if tab === "curl"}
        <div
          style="background: var(--surface); border: 1px solid var(--border);
                 border-radius: 10px; overflow: hidden;"
        >
          <div
            class="flex items-center justify-between"
            style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
          >
            <div class="flex items-center gap-2">
              <span class="material-symbols-outlined" style="font-size: 14px; color: var(--text-3);">terminal</span>
              <span
                class="uppercase font-semibold"
                style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
              >cURL Command</span>
            </div>
            <CopyBtn text={buildCurl(record)} label="Copy command" />
          </div>
          <div style="padding: 20px; background: var(--code-bg);">
            <pre
              class="font-mono"
              style="font-size: var(--fs-ui); color: var(--code-text); line-height: 1.7;"
            >{@html highlightShell(buildCurl(record))}</pre>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}
