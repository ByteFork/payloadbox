import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), svelte()],
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
