<script lang="ts">
export let text: string = "";
export let label: string = "";

let copied = false;
let timer: ReturnType<typeof setTimeout> | null = null;

function copy(): void {
  navigator.clipboard.writeText(text).then(() => {
    copied = true;
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => (copied = false), 1800);
  });
}
</script>

<button
  type="button"
  on:click={copy}
  class="inline-flex items-center gap-1 transition-all"
  style="padding: 4px 10px;
         border-radius: 6px;
         background: {copied ? 'var(--ok-bg)' : 'var(--surface2)'};
         color: {copied ? 'var(--ok)' : 'var(--text-3)'};
         border: 1px solid var(--border);
         font-size: var(--fs-data);
         font-weight: 500;"
>
  <span class="material-symbols-outlined" style="font-size: 13px;">
    {copied ? "check" : "content_copy"}
  </span>
  {#if label}<span>{copied ? "Copied" : label}</span>{/if}
</button>
