<div align="center" markdown="1">

![PayloadBox](./docs/assets/banner.svg)

Lightweight, self-hosted HTTP Request Inspector. Capture any HTTP request sent to it, inspect headers and body, and stream new captures live over Server-Sent Events - all in a single static binary.

[![build](https://img.shields.io/github/actions/workflow/status/ByteFork/payloadbox/go.yml?branch=main&label=build)](https://github.com/ByteFork/payloadbox/actions/workflows/go.yml)
[![UI checks](https://img.shields.io/github/actions/workflow/status/ByteFork/payloadbox/ui.yml?branch=main&label=ui%20checks)](https://github.com/ByteFork/payloadbox/actions/workflows/ui.yml)
[![Release](https://img.shields.io/github/v/release/ByteFork/payloadbox?sort=semver)](https://github.com/ByteFork/payloadbox/releases/latest)
[![License](https://img.shields.io/github/license/ByteFork/payloadbox)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ByteFork/payloadbox)](https://goreportcard.com/report/github.com/ByteFork/payloadbox)

</div>

## Features

- Captures arbitrary HTTP requests while reserving `/`, `/index.html`, and `/assets/*` for the embedded UI
- Embedded Svelte UI for browsing captures, inspecting headers/body/query data, and copying replayable cURL commands
- JSON API for listing, fetching, streaming, and clearing captured requests
- Server-Sent Events stream with a 30-second heartbeat for long-lived connections
- Bounded in-memory ring buffer; oldest records evicted when full
- Per-request body-size limit with graceful 413 (still recorded)
- Published as Distroless Static and Alpine Linux container images

## Usage

Start the server and send it a request:

```bash
~ payloadbox
{"time":"2026-04-23T19:00:00Z","level":"INFO","msg":"starting server","address":":8080"}

~ curl -X POST http://localhost:8080/webhooks/test \
    -H 'Content-Type: application/json' -d '{"event":"ping"}'
Request logged POST /webhooks/test

~ curl -s http://localhost:8080/api/v1/history | jq -r '.[-1].request.path'
/webhooks/test
```

Watch captures arrive live:

```bash
~ curl -N http://localhost:8080/api/v1/events
event: record
data: {"created_at":"2026-04-23T19:00:01Z",...}
```

### Container Image

```bash
~ docker run --rm -p 8080:8080 ghcr.io/bytefork/payloadbox:latest
```

## Configuration

Environment variables only.

| Variable | Default | Description |
|---|---|---|
| `LISTEN_ADDRESS` | `:8080` | Host and port to bind |
| `MAX_BODY_SIZE_BYTES` | `5120` | Per-request body limit; over-limit returns 413 but is still recorded |
| `MAX_RECORDS_TO_STORE` | `200` | Ring-buffer capacity |
| `LOG_HTTP_REQUESTS` | `true` | Log each capture to stdout |
| `LOG_LEVEL` | `info` | One of `debug`, `info`, `warn`, `error` |

## API

| Method | Path | Purpose |
|---|---|---|
| `ANY` | `/*` | Capture endpoint, except UI assets and built-in API routes |
| `GET` | `/` | Embedded web UI |
| `GET` | `/api/v1/history` | List records (gzip when accepted) |
| `GET` | `/api/v1/history/{id}` | Get one record by ID (gzip when accepted) |
| `DELETE` | `/api/v1/history` | Clear records |
| `GET` | `/api/v1/events` | SSE stream of new records |
| `GET` | `/api/v1/settings` | Current effective config |
| `GET` | `/version` | Build metadata |
| `GET` | `/healthz` | Liveness probe |

## Installation

### Install script

```bash
curl -fsSL https://install.bytefork.io/payloadbox | sh
```

Installs the latest release to `/usr/local/bin`. The script detects your OS and architecture, downloads the matching release asset, and verifies SHA-256 against the release's `checksums.txt`.

To pin a specific version:

```bash
curl -fsSL https://install.bytefork.io/payloadbox | sh -s -- --version v0.0.1
```

Pass `--help` for other options.

### Alternative methods

Download a binary directly from [Releases](https://github.com/ByteFork/payloadbox/releases), or install from source:

```bash
~ go install github.com/ByteFork/payloadbox@latest
```

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This repository is [MIT](LICENSE) licensed.
