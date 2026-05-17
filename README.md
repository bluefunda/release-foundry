# release-foundry

A shared CI/CD platform for bluefunda repos. Provides reusable GitHub Actions workflows and a binary that collects merged PR data and generates structured release notes.

## Reusable Workflows

All workflows are consumed via `uses: bluefunda/release-foundry/.github/workflows/<name>.yml@main`.

### `release-please.yml`

Automates versioning and CHANGELOG generation using [release-please](https://github.com/googleapis/release-please).

```yaml
release:
  uses: bluefunda/release-foundry/.github/workflows/release-please.yml@main
  with:
    release-type: node   # node | go | simple
  permissions:
    contents: write
    pull-requests: write
```

### `go-ci.yml`

Full Go CI pipeline: build, test (race detector + Codecov), lint (golangci-lint).

```yaml
ci:
  uses: bluefunda/release-foundry/.github/workflows/go-ci.yml@main
  with:
    goprivate: github.com/bluefunda/*
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `go-binary-release.yml`

GoReleaser-based binary release with macOS code signing, notarization, and Homebrew tap publishing. See [docs/macos-notarization.md](docs/macos-notarization.md) for setup.

```yaml
goreleaser:
  needs: release
  if: ${{ needs.release.outputs.release_created }}
  uses: bluefunda/release-foundry/.github/workflows/go-binary-release.yml@main
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
    HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
    MACOS_CERTIFICATE: ${{ secrets.MACOS_CERTIFICATE }}
    MACOS_CERTIFICATE_PWD: ${{ secrets.MACOS_CERTIFICATE_PWD }}
    NOTARIZATION_ISSUER_ID: ${{ secrets.NOTARIZATION_ISSUER_ID }}
    NOTARIZATION_KEY_ID: ${{ secrets.NOTARIZATION_KEY_ID }}
    NOTARIZATION_KEY: ${{ secrets.NOTARIZATION_KEY }}
```

### `docker-deploy.yml`

Builds and pushes a Docker image to GHCR after a release.

```yaml
deploy:
  needs: release
  if: ${{ needs.release.outputs.release_created }}
  uses: bluefunda/release-foundry/.github/workflows/docker-deploy.yml@main
  with:
    tag: ${{ needs.release.outputs.tag_name }}
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `github-release-notes.yml`

Generates labeled-PR release notes and populates the GitHub release body. Automatically detects the time window from the previous release.

```yaml
release-notes:
  needs: [release]
  if: ${{ needs.release.outputs.release_created }}
  uses: bluefunda/release-foundry/.github/workflows/github-release-notes.yml@main
  with:
    tag: ${{ needs.release.outputs.tag_name }}
    mode: replace   # replace (default) | append (use for CLI repos with goreleaser header)
  secrets:
    GH_PAT: ${{ secrets.GH_PAT }}
```

### `pr-title-check.yml`

Enforces [Conventional Commits](https://www.conventionalcommits.org/) format on PR titles.

```yaml
pr-title:
  uses: bluefunda/release-foundry/.github/workflows/pr-title-check.yml@main
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

## Project Structure

```
.github/workflows/          # Reusable workflows
cmd/release-foundry/        # Binary entrypoint
internal/
  config/                   # Multi-repo batch config loader
  domain/                   # Types, label rules, classification
  github/                   # GitHub REST API client
  service/                  # Fetch → filter → render orchestration
docs/
  cli-reference.md          # Full CLI docs, batch mode, output format
  macos-notarization.md     # macOS code signing setup guide
repos.example.yml           # Example batch config
```
