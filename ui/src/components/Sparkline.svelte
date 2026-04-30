<script lang="ts">
import type { RequestRecord } from "../types";
export let requests: RequestRecord[] = [];

const W = 100;
const H = 28;
const BUCKETS = 20;
const WINDOW_MS = 15 * 60 * 1000;

$: counts = (() => {
  const now = Date.now();
  const arr = Array(BUCKETS).fill(0);
  for (const r of requests) {
    // Clamp at 0 to handle backend/browser clock skew.
    const age = Math.max(0, now - new Date(r.created_at).getTime());
    if (age < WINDOW_MS) {
      const idx = Math.floor((age / WINDOW_MS) * BUCKETS);
      if (idx < BUCKETS) arr[BUCKETS - 1 - idx]++;
    }
  }
  return arr;
})();

$: max = Math.max(...counts, 1);
$: points = counts
  .map((c: number, i: number) => {
    const x = (i / (BUCKETS - 1)) * W;
    const y = H - (c / max) * (H - 4) - 2;
    return `${x},${y}`;
  })
  .join(" ");
</script>

<svg width={W} height={H} style="overflow: visible;">
  <polyline
    {points}
    fill="none"
    stroke="var(--primary)"
    stroke-width="1.5"
    stroke-linejoin="round"
    stroke-linecap="round"
    opacity="0.7"
  />
  {#each counts as c, i}
    {#if c > 0}
      <circle
        cx={(i / (BUCKETS - 1)) * W}
        cy={H - (c / max) * (H - 4) - 2}
        r="2"
        fill="var(--primary)"
        opacity="0.6"
      />
    {/if}
  {/each}
</svg>
