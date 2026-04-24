# Contributing

Thanks for helping improve PayloadBox.

## Before you start

Open an issue describing the change if it is non-trivial. Small fixes (typos, doc tweaks, obvious bug fixes) do not need one.

## Development

Requires Go 1.26+.

```bash
git clone https://github.com/ByteFork/payloadbox.git
cd payloadbox

go build -o payloadbox .
go test -race ./...
golangci-lint run ./...
```

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

## AI usage

AI assistance is welcome. Disclose it by including a `Co-Authored-By:` trailer in the commit message, or by mentioning the tool in the PR description. You are responsible for reviewing and understanding everything you submit.

## Reporting bugs and feature requests

Use the issue templates in `.github/ISSUE_TEMPLATE/`. For security concerns, please email the maintainers directly instead of opening a public issue.

## License

Contributions are licensed under the [MIT License](LICENSE). By opening a PR you agree your contributions will be distributed under these terms.
