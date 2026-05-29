# Contributing to release-foundry

Thank you for your interest in contributing! This document covers the development
workflow, coding standards, and how to add new renderers.

## Getting started

```bash
git clone https://github.com/<org>/release-foundry
cd release-foundry
go mod download
make test
```

Requirements: Go 1.21+ (see `go.mod` for the current minimum).

## Development workflow

```bash
make build      # build the binary
make test       # run all tests
make race       # run tests with race detector
make lint       # run golangci-lint
make vet        # run go vet
```

### Running against a real repo

```bash
export GITHUB_TOKEN=ghp_...
./release-foundry -owner golang -repo go -days 7 -render github-release -out ./out
```

## Commit conventions

All commits and PR titles must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add RSS renderer
fix(github): retry on 5xx responses
chore: bump golangci-lint to v2.12
```

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `perf`, `test`, `ci`, `style`, `build`, `revert`.

Breaking changes: append `!` after the type — `feat!: rename -output flag to -out`.

## Adding a renderer

A renderer transforms `domain.WeeklySummary` or `domain.BatchSummary` into a
specific output format (Markdown, JSON, HTML, etc.).

1. Create `internal/renderers/<name>.go` in package `renderers`.
2. Define a type that implements the `Renderer` interface:

```go
type myRenderer struct{}

func (myRenderer) FileExtension() string { return "md" }

func (myRenderer) Single(summary domain.WeeklySummary) string {
    // ... render a single-repo summary
}

func (myRenderer) Batch(batch domain.BatchSummary) string {
    // ... render a multi-repo summary (or iterate Single per repo)
}
```

3. Register it in an `init()` function:

```go
func init() {
    Register("my-format", myRenderer{})
}
```

4. Add table-driven tests in `internal/renderers/<name>_test.go`.

5. Update `docs/cli-reference.md` with the new renderer name and output format.

That's it — the registry and CLI help text update automatically.

## Pull request checklist

- [ ] `make test` passes (including race detector)
- [ ] `make lint` passes
- [ ] New functionality has tests
- [ ] PR title follows Conventional Commits
- [ ] Documentation updated if behavior changed

## Code style

See [CLAUDE.md](CLAUDE.md) for the full Go style guide. Key points:

- Standard library patterns preferred; small, single-method interfaces
- `context.Context` as first argument on IO-bound functions
- Errors wrapped with `%w`; sentinel errors for expected failure cases
- Table-driven tests; prefer real implementations over mocks

## Releasing

Releases are automated via Release Please. When a release PR is merged, GoReleaser
builds and publishes binaries for all supported platforms. You do not need to manually
tag or upload assets.
