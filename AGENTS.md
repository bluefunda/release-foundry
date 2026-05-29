# AGENTS.md

This is `release-foundry` — a GitHub PR-based release notes generator and
reusable CI/CD workflow library for Go projects.

## What this repo does

1. **Go binary** (`cmd/release-foundry/`) — collects merged PRs from GitHub,
   filters by label, and renders structured release notes in pluggable output
   formats (Markdown, JSON, etc.)
2. **Reusable workflows** (`.github/workflows/`) — generic GitHub Actions workflows
   for Go projects: CI, binary releases, Docker builds, release automation, PR title
   enforcement, and release notes generation.

## Running the binary

```bash
make build

# Single repo — last 7 days
./release-foundry -owner <org> -repo <name> -days 7 -render github-release -out ./out

# Single repo — since a specific release date
./release-foundry -owner <org> -repo <name> -since 2024-04-01T00:00:00Z \
  -render github-release -out ./out

# Multi-repo batch
./release-foundry -config repos.yml -days 7

# Topic discovery
./release-foundry -topic active -owner <org> -render github-release -out ./out

# Print version
./release-foundry version
```

The `github-release` renderer writes `<owner>-<repo>-github-release.md` to the
output directory.

## Generating social / blog content

Feed the rendered Markdown to an LLM:

```bash
./release-foundry -owner <org> -repo <name> -since <date> -render github-release -out ./out
claude -p "Given these release notes, write: (1) a LinkedIn post, (2) a tweet thread,
(3) a blog intro. $(cat ./out/*-github-release.md)"
```

## Pipeline position

release-foundry sits at the **end** of the pipeline:

```
Release created → github-release-notes.yml → binary runs → GitHub release body updated
```

## Internal packages

| Package | Role |
|---|---|
| `internal/config` | Multi-repo batch config loader |
| `internal/domain` | Types, label rules, PR classification |
| `internal/github` | GitHub REST API client (paginated, rate-limit aware) |
| `internal/renderers` | Pluggable output renderers + registry |
| `internal/service` | Collect → filter → render orchestration |

## Conventions

- See [CLAUDE.md](CLAUDE.md) for Go style and commit conventions
- New render formats go in `internal/renderers/` — implement the `Renderer`
  interface and register via `init()` — no changes to `main.go` required
- `repos.yml` (gitignored) configures multi-repo batch runs; `repos.example.yml`
  is the committed template
- All workflows default to `ubuntu-latest`; override with the `runner` input if
  needed
