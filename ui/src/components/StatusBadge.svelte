<script lang="ts">
import { statusInfo } from "../lib/format";
export let code: number | undefined | null = null;
export let text: string | undefined = "";
export let small: boolean = false;
$: info = statusInfo(code);
</script>

{#if !code}
  <span
    class="font-mono"
    style="color: var(--text-4); font-size: {small ? 10 : 11}px;"
  >—</span>
{:else}
  <span
    class="inline-block whitespace-nowrap font-mono font-semibold"
    style="color: {info.color};
           background: {info.bg};
           font-size: {small ? 10 : 11}px;
           padding: {small ? '2px 5px' : '2px 7px'};
           border-radius: 4px;"
  >{text && !small ? `${code} ${text}` : code}</span>
{/if}
