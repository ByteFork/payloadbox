<script lang="ts">
import { onMount } from "svelte";
import Docs from "./components/Docs.svelte";
import Inspector from "./components/Inspector.svelte";
import RequestList from "./components/RequestList.svelte";
import Settings from "./components/Settings.svelte";
import Sidebar from "./components/Sidebar.svelte";
import { api } from "./lib/api";
import type { PageType, RequestRecord } from "./types";

let requests: RequestRecord[] = [];
let selectedRequest: RequestRecord | null = null;
let activePage: PageType = "requests";
let listening: boolean = false;
let flashingIds: Set<string> = new Set();
// Follow-newest until the user manually picks a row.
let autoFollow: boolean = true;

let es: EventSource | null = null;
let healthTimer: ReturnType<typeof setInterval> | null = null;

function flash(id: string): void {
  flashingIds = new Set(flashingIds).add(id);
  setTimeout(() => {
    if (flashingIds.has(id)) {
      const next = new Set(flashingIds);
      next.delete(id);
      flashingIds = next;
    }
  }, 1200);
}

async function fetchHistory(): Promise<void> {
  try {
    const response = await api.getHistory();
    if (!response.ok) return;
    const raw: RequestRecord[] = await response.json();
    requests = raw.sort(
      (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    );
    if (requests.length > 0 && (autoFollow || !selectedRequest)) {
      selectedRequest = requests[0] ?? null;
    }
  } catch (err) {
    console.error("Failed to fetch history:", err);
  }
}

function openSSE(): void {
  if (es) return;
  const src = api.subscribe();
  src.addEventListener("record", (event: MessageEvent) => {
    const newRecord: RequestRecord = JSON.parse(event.data);
    requests = [newRecord, ...requests];
    flash(newRecord.id);
    if (autoFollow || !selectedRequest) selectedRequest = newRecord;
  });
  // On error, tear down so the next health-check tick can reopen. Don't touch
  // `listening` here — the poll is the single source of truth.
  src.onerror = () => {
    src.close();
    es = null;
  };
  es = src;
}

async function healthCheck(): Promise<void> {
  let ok = false;
  try {
    const res = await api.getHealth();
    ok = res.ok;
  } catch {
    /* backend unreachable */
  }

  const wasListening = listening;
  listening = ok;

  if (ok) {
    if (!es) openSSE();
    // Coming back after a drop: refresh history so we're current.
    if (!wasListening) fetchHistory();
  } else if (es) {
    es.close();
    es = null;
  }
}

onMount(() => {
  healthCheck();
  healthTimer = setInterval(healthCheck, 3000);

  return () => {
    if (healthTimer) clearInterval(healthTimer);
    es?.close();
    es = null;
  };
});

function handleSelect(request: RequestRecord): void {
  selectedRequest = request;
  // Clicking the newest (topmost) row re-engages live follow; any other row stops it.
  autoFollow = request.id === requests[0]?.id;
}

async function handleClear(): Promise<void> {
  try {
    await api.clearHistory();
    requests = [];
    selectedRequest = null;
    // Empty slate — resume follow-newest.
    autoFollow = true;
  } catch (err) {
    console.error("Failed to clear history:", err);
  }
}

function handleNavigate(page: PageType): void {
  activePage = page;
}
</script>

<div
  class="flex h-screen w-full overflow-hidden"
  style="background: var(--bg); color: var(--text);"
>
  <Sidebar
    {activePage}
    onNavigate={handleNavigate}
    total={requests.length}
    {listening}
  />
  {#if activePage === "requests"}
    <RequestList
      {requests}
      {selectedRequest}
      {flashingIds}
      onSelect={handleSelect}
      onClear={handleClear}
    />
    <Inspector record={selectedRequest} />
  {:else if activePage === "docs"}
    <Docs />
  {:else if activePage === "settings"}
    <Settings {listening} />
  {/if}
</div>
