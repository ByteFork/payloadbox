import { writeFileSync } from "node:fs";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

// Vite's `build.emptyOutDir: true` (the default) wipes ui/dist/ before each
// build, including the tracked .gitkeep. The Go binary's //go:embed all:dist
// requires at least one file to exist in dist on a fresh checkout, so this
// plugin re-creates the empty .gitkeep after every build. Writing a 0-byte
// file is a no-op against the tracked version, so git status stays clean.
const preserveGitkeep = {
  name: "preserve-gitkeep",
  closeBundle() {
    writeFileSync("dist/.gitkeep", "");
  },
};

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), svelte(), preserveGitkeep],
  server: {
    proxy: {
      "/api/v1/events": {
        target: "http://localhost:8080",
        changeOrigin: false,
        ws: false,
        // SSE: long-lived, no timeout, no buffering
        timeout: 0,
        proxyTimeout: 0,
        configure: (proxy) => {
          proxy.on("proxyRes", (proxyRes) => {
            // Ensure Node's http-proxy doesn't try to gzip/deflate the stream
            delete proxyRes.headers["content-length"];
            proxyRes.headers["cache-control"] = "no-cache";
            proxyRes.headers.connection = "keep-alive";
          });
        },
      },
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/version": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/healthz": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
