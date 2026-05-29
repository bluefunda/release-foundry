# release-foundry

A GitHub PR-based release notes generator and CI/CD workflow library. Collects
merged pull requests, filters by label and time window, and renders structured
release notes in multiple formats — no LLM required for the deterministic output,
AI-ready data for the generative layer.

---

## What it does

release-foundry sits at the end of your release pipeline. After a tag or release
is created, it:

1. **Collects** merged PRs from one or more GitHub repositories within a time window
2. **Filters** by label (feature, fix, performance, security, infrastructure)
3. **Renders** structured output into one or more formats (GitHub release body, JSON)
4. **Outputs** files ready to publish or feed to an LLM for blog posts / social content

```
git tag v1.2.0
  └─ release-foundry runs
       ├─ GitHub release body (Markdown, grouped by type)
       └─ release-summary.json (full structured data)
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
    uses: <org>/release-foundry/.github/workflows/go-ci.yml@main
```

Inputs: `go-version`, `build-command`, `test-command`, `golangci-lint-version`,
`run-codecov`, `runner`.

### `go-binary-release.yml` — GoReleaser binary release

```yaml
goreleaser:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: <org>/release-foundry/.github/workflows/go-binary-release.yml@main
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
  uses: <org>/release-foundry/.github/workflows/release-please.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

Outputs: `release_created`, `tag_name`.

### `github-release-notes.yml` — Populate GitHub release body

```yaml
release-notes:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: <org>/release-foundry/.github/workflows/github-release-notes.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
    release-foundry-repo: <org>/release-foundry
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `docker-deploy.yml` — Multi-arch Docker build + GHCR push

```yaml
deploy:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: <org>/release-foundry/.github/workflows/docker-deploy.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
    gitops-repo: myorg/gitops                   # optional
    gitops-compose-path: apps/myapp/compose.yaml # optional
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `pr-title-check.yml` — Conventional Commits enforcement

```yaml
pr-title:
  uses: <org>/release-foundry/.github/workflows/pr-title-check.yml@main
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
