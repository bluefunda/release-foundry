# AGENTS.md

This is `release-foundry` — the release tooling binary and CI/CD workflow library for bluefunda.

## What this repo does

1. **Go binary** (`cmd/release-foundry/`) — collects merged PRs, filters by label, renders release notes
2. **Reusable workflows** (`.github/workflows/`) — `go-binary-release.yml`, `docker-build.yml`; all other org-wide workflows now live in `bluefunda/.github`

## Running the binary

```bash
make build

# Single repo
./release-foundry -owner bluefunda -repo <name> -days 30 -render github-release -out ./out

# Multi-repo batch
./release-foundry -config repos.yml -days 7 -output batch.json
```

The `github-release` renderer writes `<repo>-github-release.md` to the output directory.

## Generating social media / blog content

Feed the rendered markdown to Claude:

```bash
./release-foundry -owner bluefunda -repo <name> -since <date> -render github-release -out ./out
claude -p "Given these release notes, write: (1) a LinkedIn post, (2) a tweet thread, (3) a blog intro. $(cat ./out/*-github-release.md)"
```

See [README.md](README.md) for the full social content workflow.

## Pipeline position

release-foundry sits at the **end** of the pipeline:

```
Release created → github-release-notes.yml → binary runs → GitHub release body updated
```

The `github-release-notes.yml` workflow (canonical in `bluefunda/.github`) checks out this repo, builds the binary, and runs it against the target repo.

## Internal packages

| Package | Role |
|---|---|
| `internal/config` | Multi-repo batch config loader |
| `internal/domain` | Types, label rules, PR classification |
| `internal/github` | GitHub REST API client |
| `internal/renderers` | Output renderers (`github-release`) |
| `internal/service` | Fetch → filter → render orchestration |

## Conventions

- See [CLAUDE.md](CLAUDE.md) for Go style and commit conventions
- New render formats go in `internal/renderers/` — implement the renderer interface and register in `cmd/release-foundry/main.go`
- `repos.yml` (gitignored) configures multi-repo batch runs; `repos.example.yml` is the committed template
