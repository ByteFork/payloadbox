<script lang="ts">
import { onMount } from "svelte";
import type { PageType } from "../types";

export let activePage: PageType = "requests";
export let onNavigate: (page: PageType) => void = () => {};
export let total: number = 0;
export let listening: boolean = true;

let isDark = false;

onMount(() => {
  const saved = localStorage.getItem("theme");
  const mq = window.matchMedia("(prefers-color-scheme: dark)");
  isDark = saved === "dark" || (!saved && mq.matches);
  applyTheme();

  const sysHandler = (e: MediaQueryListEvent) => {
    if (!localStorage.getItem("theme")) {
      isDark = e.matches;
      applyTheme();
    }
  };
  mq.addEventListener("change", sysHandler);
  return () => mq.removeEventListener("change", sysHandler);
});

function applyTheme(): void {
  document.documentElement.classList.toggle("dark", isDark);
}

function toggleTheme(): void {
  isDark = !isDark;
  localStorage.setItem("theme", isDark ? "dark" : "light");
  applyTheme();
}

const nav: { id: PageType; icon: string; label: string }[] = [
  { id: "requests", icon: "sensors", label: "Requests" },
  { id: "docs", icon: "menu_book", label: "Documentation" },
  { id: "settings", icon: "settings", label: "Configuration" },
];
</script>

<aside
  class="w-[52px] shrink-0 flex flex-col items-center gap-1 pt-2.5 pb-2.5"
  style="background: var(--sidebar-bg); border-right: 1px solid var(--sidebar-border);"
>
  <!-- Logo -->
  <div
    class="w-full flex justify-center pb-2.5 mb-1"
    style="border-bottom: 1px solid var(--sidebar-border);"
  >
    <button
      type="button"
      title="PayloadBox - go to Requests"
      aria-label="Go to Requests"
      on:click={() => onNavigate("requests")}
      class="w-[30px] h-[30px] flex items-center justify-center shrink-0"
      style="background: var(--primary); box-shadow: 0 0 10px var(--primary-glow); border-radius: 8px;"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18" height="18" viewBox="0 0 24 24"
        fill="none" stroke="white" stroke-width="2"
        stroke-linecap="round" stroke-linejoin="round"
        aria-hidden="true"
      >
        <path d="M19.07 4.93A10 10 0 0 0 6.99 3.34" />
        <path d="M4 6h.01" />
        <path d="M2.29 9.62A10 10 0 1 0 21.31 8.35" />
        <path d="M16.24 7.76A6 6 0 1 0 8.23 16.67" />
        <path d="M12 18h.01" />
        <path d="M17.99 11.66A6 6 0 0 1 15.77 16.67" />
        <circle cx="12" cy="12" r="2" />
        <path d="m13.41 10.59 5.66-5.66" />
      </svg>
    </button>
  </div>

  <!-- Nav -->
  <nav class="flex-1 flex flex-col items-center gap-0.5 w-full px-1.5">
    {#each nav as item (item.id)}
      {@const selected = activePage === item.id}
      <div class="relative w-full">
        <button
          type="button"
          title={item.label}
          on:click={() => onNavigate(item.id)}
          class="w-full flex items-center justify-center py-2 rounded-lg transition-colors"
          style={selected
            ? "background: var(--sidebar-active); color: var(--sidebar-text);"
            : "background: transparent; color: var(--text-4);"}
          on:mouseenter={(e) => {
            if (!selected) (e.currentTarget as HTMLElement).style.color = "var(--text-3)";
          }}
          on:mouseleave={(e) => {
            if (!selected) (e.currentTarget as HTMLElement).style.color = "var(--text-4)";
          }}
        >
          <span
            class="material-symbols-outlined"
            style="font-size: 18px; font-variation-settings: 'FILL' {selected ? 1 : 0}, 'wght' 400;"
          >{item.icon}</span>
        </button>
        {#if item.id === "requests" && total > 0}
          <span
            class="absolute top-0.5 right-0.5 font-bold rounded-md leading-[1.4] pointer-events-none"
            style="background: var(--text-4); color: #fff; font-size: 9px; padding: 1px 4px;"
          >
            {total > 99 ? "99+" : total}
          </span>
        {/if}
      </div>
    {/each}
  </nav>

  <!-- Live dot -->
  <div
    class="relative mb-1.5"
    style="width: 8px; height: 8px;"
    title={listening ? "Combobulating... on localhost" : "Offline"}
  >
    <div
      class="absolute inset-0 rounded-full"
      style="background: {listening ? '#22c55e' : '#ef4444'}; animation: pulse-dot 2s ease-in-out infinite;"
    ></div>
    <div
      class="absolute inset-0 rounded-full"
      style="background: {listening ? '#22c55e' : '#ef4444'}; animation: ping 1.5s ease-out infinite;"
    ></div>
  </div>

  <!-- Theme toggle -->
  <button
    type="button"
    title={isDark ? "Light mode" : "Dark mode"}
    on:click={toggleTheme}
    class="flex items-center justify-center rounded-lg transition-colors"
    style="width: 36px; height: 36px; color: var(--text-4);"
    on:mouseenter={(e) => ((e.currentTarget as HTMLElement).style.background = "var(--surface2)")}
    on:mouseleave={(e) => ((e.currentTarget as HTMLElement).style.background = "transparent")}
  >
    <span class="material-symbols-outlined" style="font-size: 17px;">
      {isDark ? "light_mode" : "dark_mode"}
    </span>
  </button>
</aside>
