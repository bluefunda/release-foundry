# CLAUDE.md

release-foundry — a GitHub PR-based release notes generator and reusable CI/CD
workflow library. See [AGENTS.md](AGENTS.md) for pipeline context and binary usage.

## Go style

- Standard library patterns preferred: accept interfaces, return concrete types
- Small, single-method interfaces over large "god interfaces"
- `context.Context` as first argument on all IO-bound functions
- Errors wrapped with `%w`; sentinel errors for expected failure cases
- Table-driven tests; prefer real implementations over mocks

## Conventions

- Conventional Commits for all commit messages and PR titles (`feat:`, `fix:`, `chore:`, etc.)
- Release Please manages versioning and CHANGELOG — do not manually edit CHANGELOG.md
- Use `/go-review` to run a full idiomatic Go audit of this repo

## CI

- `go test -race ./...` must pass before merging
- `golangci-lint` runs on every PR
