# release-foundry

release-foundry provides two things:

1. **Reusable GitHub Actions workflows** — a workflow library that any repo can
   call to build and publish a multi-arch Docker image to
   `ghcr.io/bluefunda/<repo>`, run CI, enforce Conventional Commits, manage
   releases, and more.

2. **A Go binary** — collects merged pull requests from one or more GitHub
   repositories, filters by label, and renders structured output (JSON + Markdown)
   ready to feed to an LLM to produce social media posts, blog intros, and
   newsletter copy.

---

## What it does

### 1. Docker image builds (reusable workflows)

Repos in the `bluefunda` org call release-foundry's reusable workflows from their
own `workflow.yml` files. The `docker-deploy.yml` workflow handles the full
multi-arch build and push to `ghcr.io/bluefunda/<repo>` — no per-repo
Dockerfile boilerplate required beyond the Dockerfile itself.

### 2. PR feed for social media content (binary)

After a release, the `release-foundry` binary collects merged PRs, filters them
by label, and writes structured JSON and Markdown files. Feed those files to an
LLM to generate LinkedIn posts, tweet threads, blog intros, or newsletters.

```
git tag v1.2.0
  └─ release-foundry binary runs
       ├─ GitHub release body (Markdown, grouped by label)
       └─ release-summary.json (structured feed for LLM consumption)
```

---

## Installation

### Homebrew (macOS / Linux)

```bash
brew install <your-tap>/release-foundry
```

See [docs/homebrew.md](docs/homebrew.md) for tap setup.

### Go install

```bash
go install github.com/release-foundry/cmd/release-foundry@latest
```

### Download binary

Download a pre-built binary from the [Releases](../../releases) page.

```bash
# macOS arm64
curl -L https://github.com/<org>/release-foundry/releases/latest/download/release-foundry_macOS_arm64.tar.gz \
  | tar xz && mv release-foundry /usr/local/bin/
```

### Linux (apt / rpm)

```bash
# Debian/Ubuntu
curl -L https://github.com/<org>/release-foundry/releases/latest/download/release-foundry_linux_amd64.deb \
  -o release-foundry.deb && sudo dpkg -i release-foundry.deb

# Fedora/RHEL
curl -L https://github.com/<org>/release-foundry/releases/latest/download/release-foundry_linux_amd64.rpm \
  -o release-foundry.rpm && sudo rpm -i release-foundry.rpm
```

---

## Quick start

```bash
# Set your GitHub token (or use the gh CLI — it's auto-detected)
export GITHUB_TOKEN=ghp_...

# Single repo, last 7 days
release-foundry -owner myorg -repo myrepo -render github-release -out ./out

# Single repo since a specific date
release-foundry -owner myorg -repo myrepo -since 2024-04-01T00:00:00Z \
  -render github-release -out ./out

# Batch mode — multiple repos from a config file
release-foundry -config repos.yml -days 14 -render github-release -out ./out

# Topic discovery — auto-discover repos tagged "active" in an org
release-foundry -topic active -owner myorg -render github-release -out ./out

# Print version
release-foundry version
```

---

## Configuration

### Single-repo mode (flags)

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `-token` | `GITHUB_TOKEN` | gh CLI / prompt | GitHub personal access token |
| `-owner` | `GITHUB_OWNER` | prompt | Repository owner or org |
| `-repo` | `GITHUB_REPO` | prompt | Repository name |
| `-days` | | `7` | Number of days to look back |
| `-since` | | | RFC3339 timestamp (overrides `-days`) |
| `-output` | | `release-summary.json` | Output JSON file path |
| `-render` | | | Comma-separated renderer names |
| `-out` | | `.` | Output directory for rendered files |
| `-config` | | | YAML config file for batch mode |
| `-topic` | | | GitHub topic to auto-discover repos |

### Batch mode (repos.yml)

Create a `repos.yml` (based on [repos.example.yml](repos.example.yml)):

```yaml
defaults:
  owner: myorg
  baseBranch: main

repos:
  - repo: api-server
    edition: enterprise
  - repo: web-ui
    edition: free
  - repo: cli-tool
    includeLabels: [feature, fix, plugin]
    excludeLabels: [internal]
```

Run:

```bash
release-foundry -config repos.yml -days 7 -render github-release -out ./out
```

### Label filtering

release-foundry collects PRs with these labels by default:

| Label | Rendered as |
|-------|-------------|
| `feature` | ✨ Features |
| `fix` | 🐛 Bug Fixes |
| `performance` | ⚡ Performance |
| `security` | 🔒 Security |
| `infrastructure` | 🏗️ Infrastructure |

PRs without a matching label fall back to conventional commit prefix inference
(`feat:` → feature, `fix:` → fix, etc.). PRs labeled `internal`, `refactor`,
or `chore` are excluded by default.

Override per-repo in `repos.yml` via `includeLabels` / `excludeLabels`.

---

## PR template integration

release-foundry extracts structured fields from PR bodies. Add to your
`.github/PULL_REQUEST_TEMPLATE.md`:

```markdown
## Customer Impact
<!-- Who benefits and how? Required for feature / perf / security PRs. -->

## Marketing Notes
<!-- Optional. Anything noteworthy for release notes, blog posts, or customer comms? -->
```

The `marketing_notes` field is included in the JSON output and can be used by
downstream LLM renderers to produce better-quality social content.

---

## Generating social content from release notes

The `github-release` renderer produces structured Markdown. Feed that to an LLM
to generate blog posts, social media posts, or newsletters:

```bash
# 1. Generate structured release notes
release-foundry -owner myorg -repo myrepo -since 2024-04-01T00:00:00Z \
  -render github-release -out ./out

# 2. Feed to Claude for social content
claude -p "$(cat <<'EOF'
You are a developer relations writer. Given these release notes, produce:
1. A LinkedIn post (3-4 sentences, professional tone)
2. A tweet thread (3 tweets, punchy, with hashtags)
3. A short blog intro paragraph (5-6 sentences)

Release notes:
$(cat ./out/*-github-release.md)
EOF
)"
```

Tips:

- **Label PRs consistently** — well-labeled PRs produce more structured notes
- **Fill in Marketing Notes** — that context flows into the JSON output and gives
  the LLM better material
- **Batch across repos** — use `-config repos.yml` for a cross-repo sprint summary

---

## Reusable GitHub Actions workflows

release-foundry ships generic, reusable workflows for Go projects.

### `go-ci.yml` — Build, test, lint

```yaml
# .github/workflows/ci.yml
on: [pull_request, push]
jobs:
  ci:
    uses: bluefunda/release-foundry/.github/workflows/go-ci.yml@main
```

Inputs: `go-version`, `build-command`, `test-command`, `golangci-lint-version`,
`run-codecov`, `runner`.

### `go-binary-release.yml` — GoReleaser binary release

```yaml
goreleaser:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/release-foundry/.github/workflows/go-binary-release.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
    HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

Supports macOS code signing and notarization. See [docs/macos-notarization.md](docs/macos-notarization.md).

### `release-please.yml` — Automated versioning and CHANGELOG

```yaml
release-please:
  uses: bluefunda/release-foundry/.github/workflows/release-please.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

Outputs: `release_created`, `tag_name`.

### `github-release-notes.yml` — Populate GitHub release body

```yaml
release-notes:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/release-foundry/.github/workflows/github-release-notes.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
    release-foundry-repo: bluefunda/release-foundry
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `docker-deploy.yml` — Multi-arch Docker build + GHCR push

Builds a `linux/amd64` + `linux/arm64` image and pushes to
`ghcr.io/<org>/<repo>:<tag>` and `ghcr.io/<org>/<repo>:latest`.

```yaml
deploy:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/release-foundry/.github/workflows/docker-deploy.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}   # only needed for private Go module dependencies
```

Inputs: `tag` (required), `image-name` (defaults to repo name), `build-args`, `runner`.

### `pr-title-check.yml` — Conventional Commits enforcement

```yaml
pr-title:
  uses: bluefunda/release-foundry/.github/workflows/pr-title-check.yml@main
```

---

## Renderer architecture

Renderers are registered via `init()` and selected at runtime by name. Adding
a new renderer requires no changes to the CLI — only a new file in
`internal/renderers/` that calls `Register("name", renderer)`.

```
release-foundry -render github-release,my-custom-renderer -out ./out
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to write a renderer.

---

## Project structure

```
cmd/release-foundry/        Binary entrypoint, flag parsing, mode dispatch
internal/
  config/                   Multi-repo batch config loader (repos.yml)
  domain/                   Core types, label rules, PR classification
  github/                   GitHub REST API client (paginated, rate-limit aware)
  renderers/                Pluggable output renderers + registry
  service/                  Collect → filter → render orchestration
.github/workflows/          Reusable GitHub Actions workflows
docs/
  cli-reference.md          Full CLI flags, batch config, output schema
  macos-notarization.md     macOS code signing + notarization setup
  homebrew.md               Homebrew tap setup guide
repos.example.yml           Example batch config
.goreleaser.yml             GoReleaser config for binary distribution
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). All PRs require:

- `make test` passing (race detector enabled)
- `make lint` passing
- PR title following [Conventional Commits](https://www.conventionalcommits.org/)

## License

Apache 2.0 — see [LICENSE](LICENSE).
