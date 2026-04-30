<script lang="ts">
import { onMount } from "svelte";
import { api } from "../lib/api";
import { formatBytes } from "../lib/format";
import { versionStore } from "../lib/stores.svelte";

interface ServerSettings {
  address: string;
  max_body_size_bytes: number;
  max_records_to_store: number;
  log_requests: boolean;
  log_level: string;
}

export let listening: boolean = false;

let settings: ServerSettings | null = null;
let loading = true;
let error: string | null = null;

onMount(async () => {
  try {
    const res = await api.getSettings();
    if (res.ok) settings = await res.json();
  } catch (err) {
    error = (err as Error).message;
    console.error("Failed to fetch settings:", err);
  } finally {
    loading = false;
  }
});

type Row = {
  key: string;
  value: string | boolean;
  desc?: string;
  status?: boolean;
};

$: sections = (() => {
  const out: { label: string; rows: Row[] }[] = [];
  out.push({
    label: "Version",
    rows: [
      { key: "version", value: versionStore.version },
      { key: "build_sha", value: versionStore.build_sha },
      { key: "build_time", value: versionStore.build_time },
    ],
  });
  if (settings) {
    out.push({
      label: "Configuration",
      rows: [
        { key: "LISTEN_ADDRESS", value: settings.address, desc: "Host and port to bind" },
        {
          key: "MAX_BODY_SIZE_BYTES",
          value: formatBytes(settings.max_body_size_bytes),
          desc: "Per-request body limit; over-limit returns 413 but is still recorded",
        },
        {
          key: "MAX_RECORDS_TO_STORE",
          value: String(settings.max_records_to_store),
          desc: "Ring-buffer capacity",
        },
        {
          key: "LOG_HTTP_REQUESTS",
          value: settings.log_requests,
          desc: "Log each capture to stdout",
        },
        {
          key: "LOG_LEVEL",
          value: settings.log_level,
          desc: "One of debug, info, warn, error",
        },
      ],
    });
  }
  out.push({
    label: "Runtime",
    rows: [
      { key: "status", value: listening ? "Combobulating" : "Disconnected", status: true },
      {
        key: "endpoint",
        value: settings?.address
          ? `localhost${settings.address.startsWith(":") ? settings.address : `:${settings.address}`}`
          : "localhost:8080",
      },
    ],
  });
  return out;
})();

$: statusColor = listening ? "#22c55e" : "#ef4444";
</script>

<main
  class="flex-1 flex flex-col overflow-hidden"
  style="background: var(--bg);"
>
  <div
    class="flex items-center gap-2.5"
    style="background: var(--surface); border-bottom: 1px solid var(--border); padding: 14px 24px;"
  >
    <span class="material-symbols-outlined" style="font-size: 16px; color: var(--text-3);">settings</span>
    <h1 style="font-size: var(--fs-page); font-weight: 700; color: var(--text); letter-spacing: -0.01em;">Configuration</h1>
  </div>

  <div class="flex-1 overflow-y-auto" style="padding: 24px;">
    {#if loading}
      <div class="flex items-center justify-center" style="padding: 48px 0;">
        <span
          class="material-symbols-outlined animate-spin"
          style="font-size: 32px; color: var(--text-4);"
        >progress_activity</span>
      </div>
    {:else if error}
      <div
        class="flex flex-col items-center gap-2"
        style="padding: 48px 24px; color: var(--err);"
      >
        <span class="material-symbols-outlined" style="font-size: 36px;">sensors_off</span>
        <p style="font-size: var(--fs-emph); font-weight: 500;">Failed to fetch settings</p>
        <p style="font-size: var(--fs-data);">{error}</p>
      </div>
    {:else}
      <div class="flex flex-col gap-6" style="max-width: 580px;">
        {#each sections as sec}
          <div>
            <div
              class="uppercase"
              style="font-size: var(--fs-micro); font-weight: 700; color: var(--text-4);
                     letter-spacing: 0.08em; margin-bottom: 6px; padding-left: 2px;"
            >{sec.label}</div>
            {#if sec.label === "Configuration"}
              <div
                style="font-size: var(--fs-micro); color: var(--text-4); margin-bottom: 8px; padding-left: 2px;"
              >
                Set the
                <code
                  class="font-mono"
                  style="background: var(--surface2); color: var(--text-3); padding: 1px 5px; border-radius: 4px;"
                >env var</code>
                before launching to override the default.
              </div>
            {/if}
            <div>
              {#each sec.rows as row, i}
                {@const isFirst = i === 0}
                {@const isLast = i === sec.rows.length - 1}
                <div
                  class="flex items-center gap-4"
                  style="background: {i % 2 === 0 ? 'var(--surface)' : 'var(--surface2)'};
                         padding: 10px 14px;
                         border: 1px solid var(--border);
                         border-top: {isFirst ? '1px solid var(--border)' : 'none'};
                         border-radius: {isFirst && isLast
                           ? '8px'
                           : isFirst
                             ? '8px 8px 0 0'
                             : isLast
                               ? '0 0 8px 8px'
                               : '0'};"
                >
                  <span
                    class="font-mono shrink-0"
                    style="font-size: var(--fs-data); color: var(--text-4); width: 200px;"
                  >{row.key}</span>

                  {#if row.status}
                    <span
                      class="font-mono font-semibold flex items-center shrink-0"
                      style="font-size: var(--fs-data); color: {statusColor}; min-width: 110px;"
                    >{row.value}<span class="dots"><span>.</span><span>.</span><span>.</span></span></span>
                  {:else if typeof row.value === "boolean"}
                    <span
                      class="font-mono font-semibold shrink-0"
                      style="font-size: var(--fs-data); color: {row.value ? 'var(--ok)' : 'var(--text-4)'}; min-width: 110px;"
                    >{row.value ? "true" : "false"}</span>
                  {:else}
                    <span
                      class="font-mono shrink-0"
                      style="font-size: var(--fs-data); color: var(--text-2); min-width: 110px;"
                    >{row.value}</span>
                  {/if}

                  {#if row.desc}
                    <span
                      class="ml-auto"
                      style="font-size: var(--fs-micro); color: var(--text-4); text-align: right; line-height: 1.5;"
                    >{row.desc}</span>
                  {/if}
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</main>

<style>
  @keyframes dot-pulse {
    0%, 60%, 100% { opacity: 0.15; }
    30%           { opacity: 1; }
  }
  .dots > span {
    animation-name: dot-pulse;
    animation-duration: 1.2s;
    animation-timing-function: ease-in-out;
    animation-iteration-count: infinite;
  }
  .dots > span:nth-child(2) { animation-delay: 200ms; }
  .dots > span:nth-child(3) { animation-delay: 400ms; }
</style>
