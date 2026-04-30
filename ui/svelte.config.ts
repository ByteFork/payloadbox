import type { SvelteConfig } from "@sveltejs/vite-plugin-svelte";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

// Consult https://svelte.dev/docs#compile-time-svelte-preprocess
// for more information about preprocessors
const config: SvelteConfig = {
  preprocess: vitePreprocess(),
};

export default config;
