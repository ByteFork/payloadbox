<script lang="ts">
import { formatDuration, formatRelTime, methodInfo } from "../lib/format";
import type { MethodFilter, RequestRecord } from "../types";
import MethodBadge from "./MethodBadge.svelte";
import Sparkline from "./Sparkline.svelte";
import StatusBadge from "./StatusBadge.svelte";

export let requests: RequestRecord[] = [];
export let selectedRequest: RequestRecord | null = null;
export let onSelect: (r: RequestRecord) => void = () => {};
export let onClear: () => void = () => {};
export let flashingIds: Set<string> = new Set();

let search = "";
let methodFilter: MethodFilter = "ALL";
const methods: MethodFilter[] = ["ALL", "GET", "POST", "PUT", "PATCH", "DELETE"];

$: filtered = requests.filter((r) => {
  const q = search.toLowerCase();
  const m = r.request.method.toUpperCase();
  const matchSearch = !q || r.request.path.toLowerCase().includes(q) || m.toLowerCase().includes(q);
  const matchMethod = methodFilter === "ALL" || m === methodFilter;
  return matchSearch && matchMethod;
});
</script>

<div
  class="flex flex-col shrink-0"
  style="width: 320px; background: var(--surface); border-right: 1px solid var(--border);"
>
  <!-- Header -->
  <div style="padding: 14px 16px 10px; border-bottom: 1px solid var(--border);">
    <div class="flex items-center justify-between" style="margin-bottom: 10px;">
      <span
        class="font-semibold uppercase"
        style="font-size: var(--fs-ui); color: var(--text-3); letter-spacing: 0.06em;"
      >Requests</span>
      <div class="flex gap-1.5">
        <button
          type="button"
          on:click={onClear}
          class="clear-btn flex items-center gap-1"
        >
          <span class="material-symbols-outlined" style="font-size: 12px;">delete_sweep</span>
          Clear
        </button>
      </div>
    </div>

    <!-- Sparkline -->
    <div class="flex items-end gap-2.5" style="margin-bottom: 10px;">
      <Sparkline {requests} />
      <div style="font-size: var(--fs-micro); color: var(--text-4); white-space: nowrap; padding-bottom: 2px;">
        last 15m
      </div>
    </div>

    <!-- Search -->
    <div
      class="flex items-center gap-2"
      style="background: var(--surface2); border: 1px solid var(--border);
             border-radius: 8px; padding: 6px 10px; margin-bottom: 8px;"
    >
      <span class="material-symbols-outlined" style="font-size: 14px; color: var(--text-4);">search</span>
      <input
        type="text"
        bind:value={search}
        placeholder="Filter requests…"
        class="flex-1"
        style="font-size: var(--fs-ui); color: var(--text); background: transparent;"
      />
      {#if search}
        <button type="button" on:click={() => (search = "")} style="color: var(--text-4);">
          <span class="material-symbols-outlined" style="font-size: 13px;">close</span>
        </button>
      {/if}
    </div>

    <!-- Method filter -->
    <div class="flex gap-1 flex-wrap">
      {#each methods as m}
        {@const info = m === "ALL" ? { color: "#fff", bg: "var(--primary)" } : methodInfo(m)}
        {@const active = methodFilter === m}
        <button
          type="button"
          on:click={() => (methodFilter = m)}
          class="font-mono font-semibold transition-all"
          style="padding: 3px 8px; border-radius: 5px; font-size: 10px; letter-spacing: 0.03em;
                 background: {active ? (m === 'ALL' ? 'var(--primary)' : info.bg) : 'var(--surface2)'};
                 color: {active ? (m === 'ALL' ? '#fff' : info.color) : 'var(--text-4)'};
                 border: 1px solid {active ? 'transparent' : 'var(--border)'};"
        >{m}</button>
      {/each}
    </div>
  </div>

  <!-- List -->
  <div class="flex-1 overflow-y-auto">
    {#if filtered.length === 0}
      <div class="flex flex-col items-center gap-2.5" style="padding: 48px 24px;">
        <span
          class="material-symbols-outlined"
          style="font-size: 36px; color: var(--text-4); opacity: 0.5;"
        >inbox</span>
        <span style="font-size: var(--fs-ui); color: var(--text-4);">
          No requests {search || methodFilter !== "ALL" ? "match your filter" : "yet"}
        </span>
      </div>
    {:else}
      {#each filtered as r (r.id)}
        {@const selected = selectedRequest?.id === r.id}
        {@const status = r.response?.status_code}
        {@const method = (r.request?.method ?? "").toUpperCase()}
        <div
          role="button"
          tabindex="0"
          class="req-row"
          class:selected
          class:req-row-new={flashingIds.has(r.id)}
          on:click={() => onSelect(r)}
          on:keydown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              onSelect(r);
            }
          }}
        >
          <div class="flex items-center gap-2" style="margin-bottom: 5px;">
            <MethodBadge {method} small />
            {#if status}<StatusBadge code={status} small />{/if}
            <span
              class="ml-auto font-mono"
              style="font-size: var(--fs-micro); color: var(--text-4);"
            >{formatDuration(r.duration_ns)}</span>
          </div>
          <div
            class="font-mono truncate"
            style="font-size: var(--fs-data); color: var(--text-2); margin-bottom: 3px;"
          >
            {r.request?.path ?? ""}{#if r.request?.query}<span style="color: var(--text-4);">?{r.request.query}</span>{/if}
          </div>
          <div class="flex items-center gap-1.5">
            <span style="font-size: var(--fs-micro); color: var(--text-4);">{formatRelTime(r.created_at)}</span>
            <span style="font-size: var(--fs-micro); color: var(--text-4);">·</span>
            <span class="font-mono" style="font-size: var(--fs-micro); color: var(--text-4);">{r.request?.remote_addr ?? ""}</span>
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .req-row {
    padding: 10px 16px;
    cursor: pointer;
    border-bottom: 1px solid var(--border);
    border-left: 2px solid transparent;
    background: transparent;
    transition: background 0.1s;
  }
  .req-row:hover {
    background: var(--surface2);
  }
  .req-row.selected {
    background: var(--surface2);
    border-left-color: var(--border-strong);
  }

  .clear-btn {
    padding: 4px 8px;
    border-radius: 6px;
    font-size: var(--fs-data);
    font-weight: 500;
    color: var(--text-3);
    background: var(--surface2);
    border: 1px solid var(--border);
    transition: color 0.1s, border-color 0.1s;
  }
  .clear-btn:hover {
    color: var(--err);
    border-color: var(--err);
  }
</style>
