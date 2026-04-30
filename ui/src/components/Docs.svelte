<script lang="ts">
import { highlightJSON, highlightShell, methodInfo } from "../lib/format";
import { versionStore } from "../lib/stores.svelte";
import MethodBadge from "./MethodBadge.svelte";

type DocSection = {
  id: string;
  group: string;
  label: string;
  icon: string;
  method?: string;
  path?: string;
};

const sections: DocSection[] = [
  { id: "overview", group: "Getting Started", label: "What is PayloadBox?", icon: "info" },
  { id: "quickstart", group: "Getting Started", label: "Quick Start", icon: "rocket_launch" },
  {
    id: "list",
    group: "API Reference",
    label: "List Requests",
    icon: "list",
    method: "GET",
    path: "/api/v1/history",
  },
  {
    id: "get",
    group: "API Reference",
    label: "Get Request",
    icon: "search",
    method: "GET",
    path: "/api/v1/history/{id}",
  },
  {
    id: "delete",
    group: "API Reference",
    label: "Clear Requests",
    icon: "delete_sweep",
    method: "DELETE",
    path: "/api/v1/history",
  },
  {
    id: "settings-api",
    group: "API Reference",
    label: "Settings",
    icon: "settings",
    method: "GET",
    path: "/api/v1/settings",
  },
];

let active = "overview";
$: section = sections.find((s) => s.id === active) ?? sections[0];
$: groups = Array.from(new Set(sections.map((s) => s.group)));

const dockerCmd = `docker run --rm -p 8080:8080 ghcr.io/bytefork/payloadbox:latest`;
const binaryCmd = `curl -fsSL https://install.bytefork.io/payloadbox | sh`;

const CODE = {
  list: `[
  {
    "id": "0190d6f1-9c2a-7c9d-8e3a-1f0a8b3c4d5e",
    "request": { "method": "POST", "path": "/webhooks/payment" },
    "response": { "status_code": 200, "size_in_bytes": 24 },
    "duration_ns": 4820000,
    "created_at": "2026-04-21T14:16:28Z"
  }
]`,
  get: `{
  "id": "0190d6f1-9c2a-7c9d-8e3a-1f0a8b3c4d5e",
  "created_at": "2026-04-21T14:16:28Z",
  "duration_ns": 4820000,
  "request": {
    "method": "POST",
    "path": "/webhooks/payment",
    "query": "",
    "headers": { "Content-Type": ["application/json"] },
    "body": "{\\"type\\":\\"payment.created\\"}",
    "remote_addr": "192.168.1.45",
    "host": "localhost:8080",
    "content_length": 38
  },
  "response": {
    "status_code": 200,
    "status_text": "OK",
    "headers": { "Content-Type": ["application/json"] },
    "body": "{\\"received\\":true}",
    "size_in_bytes": 17
  }
}`,
  settings: `{
  "address": ":8080",
  "max_body_size_bytes": 5242880,
  "max_records_to_store": 1000,
  "log_requests": true,
  "log_level": "info"
}`,
};
</script>

<div
  class="flex-1 flex flex-col overflow-hidden"
  style="background: var(--surface);"
>
  <!-- Header -->
  <div
    class="flex items-center gap-2.5"
    style="background: var(--surface); border-bottom: 1px solid var(--border); padding: 14px 24px;"
  >
    <span class="material-symbols-outlined" style="font-size: 16px; color: var(--primary);">menu_book</span>
    <h1 style="font-size: var(--fs-page); font-weight: 700; color: var(--text); letter-spacing: -0.01em;">Documentation</h1>
    <span class="font-mono" style="font-size: var(--fs-data); color: var(--text-4); margin-left: 4px;">
      v{versionStore.version}
    </span>
  </div>

  <div class="docs-body flex-1 flex overflow-hidden">
    <!-- Left nav -->
    <div
      class="docs-nav shrink-0 overflow-y-auto"
      style="width: 210px; border-right: 1px solid var(--border); padding: 16px 10px;"
    >
      {#each groups as group}
        <div style="margin-bottom: 18px;">
          <div
            class="uppercase"
            style="font-size: var(--fs-micro); font-weight: 700; color: var(--text-4);
                   letter-spacing: 0.08em; padding: 0 8px; margin-bottom: 4px;"
          >{group}</div>
          {#each sections.filter((s) => s.group === group) as s}
            {@const selected = active === s.id}
            <button
              type="button"
              class="w-full flex items-center gap-2 text-left transition-all"
              on:click={() => (active = s.id)}
              style="padding: 7px 8px; border-radius: 7px; margin-bottom: 1px;
                     background: {selected ? 'var(--sidebar-active)' : 'transparent'};
                     color: {selected ? 'var(--text)' : 'var(--text-3)'};
                     font-weight: {selected ? 600 : 400}; font-size: var(--fs-ui);"
              on:mouseenter={(e) => {
                if (!selected) (e.currentTarget as HTMLElement).style.background = "var(--surface2)";
              }}
              on:mouseleave={(e) => {
                if (!selected) (e.currentTarget as HTMLElement).style.background = "transparent";
              }}
            >
              {#if s.method}
                {@const info = methodInfo(s.method)}
                <span
                  class="font-mono shrink-0"
                  style="font-size: 8px; font-weight: 700;
                         color: {info.color}; background: {info.bg};
                         padding: 1px 4px; border-radius: 3px;"
                >{s.method}</span>
              {:else}
                <span class="material-symbols-outlined shrink-0" style="font-size: 13px;">{s.icon}</span>
              {/if}
              <span class="truncate">{s.label}</span>
            </button>
          {/each}
        </div>
      {/each}
    </div>

    <!-- Content -->
    <div class="docs-content flex-1 overflow-y-auto" style="padding: 32px; background: var(--surface);">
      <div style="max-width: 640px;">
        <!-- Title block -->
        <div style="margin-bottom: 24px;">
          <div class="flex items-center gap-2" style="margin-bottom: 6px;">
            <span style="font-size: var(--fs-data); color: var(--text-4);">{section.group}</span>
            <span class="material-symbols-outlined" style="font-size: 12px; color: var(--text-4);">chevron_right</span>
            {#if section.method}<MethodBadge method={section.method} />{/if}
          </div>
          <h2 style="font-size: var(--fs-h2); font-weight: 700; color: var(--text); letter-spacing: -0.02em; margin-bottom: 4px;">
            {section.label}
          </h2>
          {#if section.path}
            <code class="font-mono" style="font-size: var(--fs-ui); color: var(--text-4);">{section.path}</code>
          {/if}
        </div>

        <div style="border-top: 1px solid var(--border); padding-top: 24px;">
          {#if section.id === "overview"}
            <p style="font-size: var(--fs-body); color: var(--text-2); line-height: 1.8; margin-bottom: 20px;">
              PayloadBox is a lightweight, self-hosted HTTP request inspector. Run it locally to capture,
              inspect, and replay incoming webhook calls and API requests in real time - without exposing
              anything to the internet.
            </p>
            <div class="grid gap-3" style="grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); margin-bottom: 24px;">
              {#each [
                { icon: "wifi",     label: "Live capture", desc: "See requests the moment they arrive" },
                { icon: "search",   label: "Deep inspect", desc: "Headers, body, query, response" },
                { icon: "terminal", label: "cURL replay",  desc: "One-click copy to replay any request" },
              ] as f}
                <div style="background: var(--surface2); border: 1px solid var(--border); border-radius: 10px; padding: 14px 16px;">
                  <div class="flex items-center gap-2" style="margin-bottom: 4px;">
                    <span class="material-symbols-outlined" style="font-size: 18px; color: var(--primary);">{f.icon}</span>
                    <div style="font-size: var(--fs-emph); font-weight: 600; color: var(--text);">{f.label}</div>
                  </div>
                  <div style="font-size: var(--fs-ui); color: var(--text-4); line-height: 1.5;">{f.desc}</div>
                </div>
              {/each}
            </div>
            <div
              class="flex gap-2.5"
              style="background: var(--warn-bg); border: 1px solid var(--warn);
                     border-radius: 10px; padding: 12px 16px;"
            >
              <span class="material-symbols-outlined shrink-0" style="font-size: 16px; color: var(--warn); margin-top: 1px;">lightbulb</span>
              <span style="font-size: var(--fs-ui); color: var(--text-2); line-height: 1.6;">
                PayloadBox stores everything in memory. Requests are lost on restart unless you export them first.
              </span>
            </div>

          {:else if section.id === "quickstart"}
            <div class="flex flex-col gap-5">
              <!-- Step 1: Install — two options side by side -->
              <div>
                <div class="flex items-center gap-2.5" style="margin-bottom: 10px;">
                  <div
                    class="flex items-center justify-center shrink-0"
                    style="width: 22px; height: 22px; border-radius: 50%;
                           background: var(--primary); color: #fff;
                           font-size: var(--fs-data); font-weight: 700;"
                  >1</div>
                  <span style="font-size: var(--fs-body); font-weight: 600; color: var(--text);">Install and run</span>
                  <span
                    class="font-mono ml-auto"
                    style="font-size: var(--fs-micro); color: var(--text-4);"
                  >pick one</span>
                </div>

                <!-- Docker -->
                <div style="margin-bottom: 12px;">
                  <div
                    class="font-mono uppercase"
                    style="font-size: var(--fs-micro); color: var(--text-4);
                           letter-spacing: 0.08em; margin-bottom: 6px;"
                  >Docker</div>
                  <div style="background: var(--code-bg); border-radius: 10px; overflow: hidden;">
                    <pre
                      class="font-mono"
                      style="font-size: var(--fs-ui); color: var(--code-text); padding: 14px 16px; line-height: 1.8; overflow-x: auto;"
                    >{@html highlightShell(dockerCmd)}</pre>
                  </div>
                </div>

                <!-- Binary -->
                <div>
                  <div
                    class="font-mono uppercase"
                    style="font-size: var(--fs-micro); color: var(--text-4);
                           letter-spacing: 0.08em; margin-bottom: 6px;"
                  >Binary (Linux & macOS)</div>
                  <div style="background: var(--code-bg); border-radius: 10px; overflow: hidden;">
                    <pre
                      class="font-mono"
                      style="font-size: var(--fs-ui); color: var(--code-text); padding: 14px 16px; line-height: 1.8; overflow-x: auto;"
                    >{@html highlightShell(binaryCmd)}</pre>
                  </div>
                </div>
              </div>

              <!-- Step 2 & 3 -->
              {#each [
                { step: "2", title: "Send a request", code: `curl -X POST http://localhost:8080/test \\\n  -H 'Content-Type: application/json' \\\n  -d '{"hello":"world"}'` },
                { step: "3", title: "Open the UI", code: `open http://localhost:8080` },
              ] as s}
                <div>
                  <div class="flex items-center gap-2.5" style="margin-bottom: 10px;">
                    <div
                      class="flex items-center justify-center shrink-0"
                      style="width: 22px; height: 22px; border-radius: 50%;
                             background: var(--primary); color: #fff;
                             font-size: var(--fs-data); font-weight: 700;"
                    >{s.step}</div>
                    <span style="font-size: var(--fs-body); font-weight: 600; color: var(--text);">{s.title}</span>
                  </div>
                  <div style="background: var(--code-bg); border-radius: 10px; overflow: hidden;">
                    <pre
                      class="font-mono"
                      style="font-size: var(--fs-ui); color: var(--code-text); padding: 14px 16px; line-height: 1.8; overflow-x: auto;"
                    >{@html highlightShell(s.code)}</pre>
                  </div>
                </div>
              {/each}
              <div
                class="flex gap-2.5"
                style="background: var(--info-bg); border: 1px solid var(--info);
                       border-radius: 10px; padding: 12px 16px;"
              >
                <span class="material-symbols-outlined shrink-0" style="font-size: 16px; color: var(--info); margin-top: 1px;">settings</span>
                <span style="font-size: var(--fs-ui); color: var(--text-2); line-height: 1.6;">
                  Runtime settings and environment variable names live in the Configuration page.
                </span>
              </div>
            </div>

          {:else if section.id === "list"}
            <div class="flex flex-col gap-4">
              <p style="font-size: var(--fs-emph); color: var(--text-3); line-height: 1.7;">
                Returns the list of all captured requests, sorted newest-first.
              </p>
              <div style="background: var(--surface); border: 1px solid var(--border); border-radius: 10px; overflow: hidden;">
                <div
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response Body</span>
                </div>
                <div style="padding: 16px;">
                  <pre
                    class="font-mono"
                    style="font-size: var(--fs-ui); line-height: 1.7; color: var(--text-2);"
                  >{@html highlightJSON(CODE.list)}</pre>
                </div>
              </div>
            </div>

          {:else if section.id === "get"}
            <div class="flex flex-col gap-4">
              <p style="font-size: var(--fs-emph); color: var(--text-3); line-height: 1.7;">
                Fetch a single captured request by its ID, including full headers, body, and response data.
              </p>
              <div style="background: var(--surface); border: 1px solid var(--border); border-radius: 10px; overflow: hidden;">
                <div
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response Body</span>
                </div>
                <div style="padding: 16px;">
                  <pre
                    class="font-mono"
                    style="font-size: var(--fs-ui); line-height: 1.7; color: var(--text-2);"
                  >{@html highlightJSON(CODE.get)}</pre>
                </div>
              </div>
            </div>

          {:else if section.id === "delete"}
            <div class="flex flex-col gap-4">
              <p style="font-size: var(--fs-emph); color: var(--text-3); line-height: 1.7;">
                Clears all stored requests from memory. This action is immediate and irreversible.
              </p>
              <div
                class="flex gap-2.5"
                style="background: var(--err-bg); border: 1px solid var(--err);
                       border-radius: 10px; padding: 12px 16px;"
              >
                <span class="material-symbols-outlined shrink-0" style="font-size: 15px; color: var(--err); margin-top: 1px;">warning</span>
                <span style="font-size: var(--fs-ui); color: var(--text-2); line-height: 1.6;">
                  All captured requests are permanently deleted. There is no undo.
                </span>
              </div>
              <div style="background: var(--surface); border: 1px solid var(--border); border-radius: 10px; overflow: hidden;">
                <div
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response</span>
                </div>
                <div
                  class="flex items-center gap-2"
                  style="padding: 14px 16px;"
                >
                  <span
                    class="inline-block whitespace-nowrap font-mono font-semibold"
                    style="color: var(--info); background: var(--info-bg);
                           font-size: var(--fs-data); padding: 2px 7px; border-radius: 4px;"
                  >204 No Content</span>
                  <span style="font-size: var(--fs-ui); color: var(--text-4);">Empty response body</span>
                </div>
              </div>
            </div>

          {:else if section.id === "settings-api"}
            <div class="flex flex-col gap-4">
              <p style="font-size: var(--fs-emph); color: var(--text-3); line-height: 1.7;">
                Returns the current server configuration as JSON.
              </p>
              <div style="background: var(--surface); border: 1px solid var(--border); border-radius: 10px; overflow: hidden;">
                <div
                  style="padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--surface2);"
                >
                  <span
                    class="uppercase font-semibold"
                    style="font-size: var(--fs-data); color: var(--text-3); letter-spacing: 0.05em;"
                  >Response Body</span>
                </div>
                <div style="padding: 16px;">
                  <pre
                    class="font-mono"
                    style="font-size: var(--fs-ui); line-height: 1.7; color: var(--text-2);"
                  >{@html highlightJSON(CODE.settings)}</pre>
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  @media (max-width: 768px) {
    .docs-body { flex-direction: column; }
    .docs-nav {
      width: 100%;
      max-height: 40vh;
      border-right: none;
      border-bottom: 1px solid var(--border);
    }
    .docs-content { padding: 20px !important; }
  }
</style>
