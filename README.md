# release-foundry

A shared CI/CD platform for bluefunda repos. Provides reusable GitHub Actions workflows and a binary that collects merged PR data and generates structured release notes.

---

## Pipeline Overview

Every bluefunda Go repo follows this end-to-end flow:

```
PR opened
  └─ go-review.yml        advisory Claude Code review posted as PR comment

PR merged to main
  └─ pr-title-check.yml   conventional commit format enforced
  └─ go-ci.yml            build, test, lint

Push to main
  └─ release-please.yml   bumps version, updates CHANGELOG, opens release PR

Release PR merged
  └─ docker-deploy.yml    builds image, pushes to GHCR, updates gitops repo
  └─ github-release-notes.yml   generates structured release notes, populates GitHub release body
```

Gitops repo is watched by [Komodo](https://komo.do) — it deploys automatically when the compose file is updated.

---

## Reusable Workflows

Org-wide workflows live in [`bluefunda/.github`](https://github.com/bluefunda/.github/tree/main/.github/workflows) and are consumed via `uses: bluefunda/.github/.github/workflows/<name>.yml@main`.

### `go-review.yml`

Advisory PR review via Claude Code. Posts a non-blocking comment with idiomatic Go feedback on the PR diff.

```yaml
on:
  pull_request:
    branches: [main]

jobs:
  go-review:
    uses: bluefunda/.github/.github/workflows/go-review.yml@main
    with:
      pr-number: ${{ github.event.pull_request.number }}
    secrets:
      GH_PAT: ${{ secrets.GH_PAT }}
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

Requires the `ANTHROPIC_API_KEY` org secret (already scoped to private repos).

### `go-ci.yml`

Full Go CI pipeline: build, test (race detector), lint (golangci-lint).

```yaml
ci:
  uses: bluefunda/.github/.github/workflows/go-ci.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
  # with:
  #   private-modules: true
  #   build-command: go build ./cmd/...
  #   test-command: go test -race -coverprofile=coverage.txt ./...
  #   run-codecov: true
```

### `release-please.yml`

Automates versioning and CHANGELOG generation. Outputs `release_created` and `tag_name` for downstream jobs.

```yaml
release-please:
  uses: bluefunda/.github/.github/workflows/release-please.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
  # with:
  #   release-type: go   # default; use 'simple' for non-module repos
```

### `docker-deploy.yml`

Builds a multi-arch Docker image, pushes to GHCR, and optionally updates the gitops repo so Komodo picks up the new tag.

```yaml
deploy:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/.github/.github/workflows/docker-deploy.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
    komodo-stack: my-app   # set to your Komodo stack name
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `github-release-notes.yml`

Runs the `release-foundry` binary, generates labeled-PR release notes, and populates the GitHub release body. Automatically detects the time window from the previous release.

```yaml
release-notes:
  needs: release-please
  if: needs.release-please.outputs.release_created == 'true'
  uses: bluefunda/.github/.github/workflows/github-release-notes.yml@main
  with:
    tag: ${{ needs.release-please.outputs.tag_name }}
    mode: replace   # replace (default) | append
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `go-binary-release.yml`

GoReleaser-based binary release for CLI tools (macOS code signing, Homebrew tap). Lives in this repo since it's Go-binary-specific.

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

See [docs/macos-notarization.md](docs/macos-notarization.md) for macOS signing setup.

### `pr-title-check.yml`

Enforces [Conventional Commits](https://www.conventionalcommits.org/) format on PR titles.

```yaml
pr-title:
  uses: bluefunda/.github/.github/workflows/pr-title-check.yml@main
```

---

## The `release-foundry` Binary

Collects merged pull requests from a GitHub repo, filters by label, and renders structured release notes.

```bash
make build

# Single repo — last 30 days
./release-foundry -owner bluefunda -repo cai-bff -days 30 -render github-release -out ./out

# Single repo — since a specific release
./release-foundry -owner bluefunda -repo cai-bff -since 2026-04-01T00:00:00Z -render github-release -out ./out

# Multi-repo batch
./release-foundry -config repos.yml -days 7 -output batch.json
```

See [docs/cli-reference.md](docs/cli-reference.md) for the full flag reference, batch config format, filtering rules, and output schema.

---

## Generating Social Media & Blog Content

The `github-release` renderer produces structured markdown release notes. Use that output as input to Claude to generate social media posts, blog drafts, or internal announcements.

### From the CLI

```bash
# 1. Generate release notes into a file
./release-foundry -owner bluefunda -repo cai-bff \
  -since 2026-04-01T00:00:00Z -render github-release -out ./out

# 2. Feed to Claude Code for social/blog content
claude -p "$(cat <<'EOF'
You are a developer relations writer. Given these release notes, produce:

1. A LinkedIn post (3-4 sentences, professional tone, highlight the top user-facing change)
2. A Twitter/X thread (3 tweets max, punchy, include relevant hashtags)
3. A short blog intro paragraph (5-6 sentences, explain what changed and why it matters)

Release notes:
$(cat ./out/*-github-release.md)
EOF
)"
```

### From a Claude Code session

Open Claude Code in the repo directory after a release and run:

```
/go-review    # review if needed

Then ask:
"Generate a LinkedIn post and a tweet thread from the latest GitHub release notes for this repo."
```

Claude Code will read the release notes directly from the GitHub release or local output files.

### Tips for better output

- **Label PRs consistently** — release-foundry groups output by label (`feature`, `fix`, `performance`, etc.). Well-labeled PRs produce more structured notes and better social content.
- **Fill in the Marketing Notes field** in the PR template — that context flows directly into release notes and gives Claude better material to work with.
- **Batch across repos** — use `-config repos.yml` to produce a cross-repo summary for a sprint or milestone post.

---

## Project Structure

```
.github/workflows/          # Reusable workflows (go-binary-release, docker-build)
cmd/release-foundry/        # Binary entrypoint
internal/
  config/                   # Multi-repo batch config loader
  domain/                   # Types, label rules, classification
  github/                   # GitHub REST API client
  renderers/                # Output renderers (github-release)
  service/                  # Fetch → filter → render orchestration
docs/
  cli-reference.md          # Full CLI docs, batch mode, output format
  macos-notarization.md     # macOS code signing setup guide
repos.example.yml           # Example batch config
```
