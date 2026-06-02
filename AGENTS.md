# AGENTS.md

This is `release-foundry` — it provides two things for repos in the
`bluefunda` org:

1. **Reusable GitHub Actions workflows** — any repo calls these from its own
   `workflow.yml` to build and publish a multi-arch Docker image to
   `ghcr.io/bluefunda/<repo>`, run Go CI, enforce Conventional Commits,
   manage releases via Release Please, and more.

2. **A Go binary** (`cmd/release-foundry/`) — collects merged PRs from one or
   more GitHub repositories (single repo, list via config file, or discovered
   by topic), filters by label, and renders structured JSON + Markdown files
   ready to feed to an LLM for social media content (LinkedIn posts, tweet
   threads, blog intros, newsletters).

## Function 1 — Docker build framework (reusable workflows)

Repos call `docker-deploy.yml` from `bluefunda/release-foundry` to build and
push a multi-arch image (`linux/amd64` + `linux/arm64`) to GHCR:

```yaml
# In your repo's .github/workflows/workflow.yml
deploy:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/release-foundry/.github/workflows/docker-deploy.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}   # only needed for private Go module deps
```

The image is published as `ghcr.io/bluefunda/<repo>:<tag>` and
`ghcr.io/bluefunda/<repo>:latest`.

### Available reusable workflows

| Workflow | Purpose |
|---|---|
| `docker-deploy.yml` | Multi-arch Docker build + push to GHCR |
| `go-ci.yml` | Build, test (race detector), lint |
| `go-binary-release.yml` | GoReleaser binary release |
| `release-please.yml` | Automated versioning + CHANGELOG |
| `github-release-notes.yml` | Populate GitHub release body |
| `pr-title-check.yml` | Conventional Commits enforcement |

## Function 2 — PR feed binary (social media content)

### Running the binary

```bash
make build

# Single repo — last 7 days
./release-foundry -owner bluefunda -repo myrepo -days 7 -render github-release -out ./out

# Single repo — since a specific release date
./release-foundry -owner bluefunda -repo myrepo -since 2024-04-01T00:00:00Z \
  -render github-release -out ./out

# Multi-repo batch
./release-foundry -config repos.yml -days 7

# Topic discovery
./release-foundry -topic active -owner bluefunda -render github-release -out ./out

# Print version
./release-foundry version
```

### Generating social / blog content

```bash
./release-foundry -owner bluefunda -repo myrepo -since <date> -render github-release -out ./out
claude -p "Given these release notes, write: (1) a LinkedIn post, (2) a tweet thread,
(3) a blog intro. $(cat ./out/*-github-release.md)"
```

### Output renderers

The `github-release` renderer writes `<owner>-<repo>-github-release.md` to the
output directory. New formats go in `internal/renderers/` — implement the
`Renderer` interface and register via `init()`. No changes to `main.go` needed.

### Pipeline position

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

## Go style

- Standard library patterns preferred: accept interfaces, return concrete types
- Small, single-method interfaces over large "god interfaces"
- `context.Context` as first argument on all IO-bound functions
- Errors wrapped with `%w`; sentinel errors for expected failure cases
- Table-driven tests; prefer real implementations over mocks

## Conventions

- Conventional Commits for all commit messages and PR titles (`feat:`, `fix:`, `chore:`, etc.)
- Release Please manages versioning and CHANGELOG — do not manually edit CHANGELOG.md
- `repos.yml` (gitignored) configures multi-repo batch runs; `repos.example.yml` is the committed template
- All workflows default to `ubuntu-latest`; override with the `runner` input if needed
- Use `/go-review` to run a full idiomatic Go audit of this repo

## CI

- `go test -race ./...` must pass before merging
- `golangci-lint` runs on every PR
