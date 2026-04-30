# Contributing

Thanks for helping improve PayloadBox.

## Before you start

Open an issue describing the change if it is non-trivial. Small fixes (typos, doc tweaks, obvious bug fixes) do not need one.

## Development

Requires Go 1.26+, Node (version pinned in `ui/.nvmrc`), and pnpm.

```bash
git clone https://github.com/ByteFork/payloadbox.git
cd payloadbox

# Build the UI bundle. The Go binary embeds ui/dist at compile time
# (//go:embed all:dist in ui/embed.go), so this step must precede go build.
pnpm --dir ui install
pnpm --dir ui build

go build -o payloadbox .
go test -race ./...
golangci-lint run ./...
```

### UI workflow

```bash
cd ui

pnpm dev          # Vite dev server; proxies /api, /version, /healthz, /api/v1/events
pnpm build        # writes ui/dist (consumed by go build)
pnpm check        # svelte-check (TypeScript + Svelte)
pnpm lint         # biome + oxlint
pnpm fix          # biome check --write && oxlint --fix
pnpm knip         # unused exports / dependencies
pnpm test:e2e     # Playwright; auto-starts ../payloadbox via playwright.config.ts
```

For UI iteration: run `./payloadbox` at the repo root in one terminal, then `pnpm dev` in `ui/` in another. Vite proxies API calls to `:8080`.

## Pull request workflow

1. Fork and create a branch: `git checkout -b type/short-description` (e.g. `fix/sse-heartbeat`, `feat/webhook-replay`).
2. Keep commits focused. [Conventional Commits](https://www.conventionalcommits.org/) are preferred but not required.
3. Run `go test -race ./...` and `golangci-lint run ./...` locally. Both must pass.
4. Open the PR against `main`. Fill in the template.
5. CI must be green and the PR needs one approval before merge.
6. PRs are squash-merged into `main`.

## Coding standards

- `gofmt` and `goimports` clean (handled by `golangci-lint fmt`).
- `golangci-lint` must pass with 0 issues (see `.golangci.yml`).
- Tests accompany new functionality. The `internal/*` packages are at or near 100% line coverage; try not to regress that.
- No `TODO`/`FIXME` without a linked issue number.
- For UI changes: `pnpm check`, `pnpm lint`, and `pnpm knip` must pass; e2e specs (`pnpm test:e2e`) cover new user-visible behavior.

## AI usage

AI assistance is welcome. Disclose it by including a `Co-Authored-By:` trailer in the commit message, or by mentioning the tool in the PR description. You are responsible for reviewing and understanding everything you submit.

## Reporting bugs and feature requests

Use the issue templates in `.github/ISSUE_TEMPLATE/`. For security concerns, please email the maintainers directly instead of opening a public issue.

## License

Contributions are licensed under the [MIT License](LICENSE). By opening a PR you agree your contributions will be distributed under these terms.
