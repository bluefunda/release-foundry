# CLI Reference

## Synopsis

```
release-foundry [flags]
release-foundry version
release-foundry help
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-token` | string | `GITHUB_TOKEN` env → `gh auth token` → prompt | GitHub personal access token |
| `-owner` | string | `GITHUB_OWNER` env → prompt | Repository owner or org |
| `-repo` | string | `GITHUB_REPO` env → prompt | Repository name |
| `-days` | int | `7` | Number of days to look back for merged PRs |
| `-since` | string | (none) | Fetch PRs merged after this RFC3339 timestamp (overrides `-days`) |
| `-output` | string | `release-summary.json` | JSON output file path |
| `-render` | string | (none) | Comma-separated renderer names (e.g. `github-release`) |
| `-out` | string | `.` | Output directory for rendered artifacts |
| `-config` | string | (none) | Path to multi-repo YAML config (batch mode) |
| `-topic` | string | (none) | GitHub topic filter for auto-discovering repos in an org |

## Operating modes

### Single-repo mode

Default mode. Processes one repository specified via flags or environment variables.

```bash
release-foundry -owner myorg -repo myrepo -days 14 -render github-release -out ./out
```

Token resolution order: `-token` flag → `GITHUB_TOKEN` env → `gh auth token` → interactive prompt.

Owner/repo resolution order: flag → environment variable → interactive prompt.

### Batch mode (`-config`)

Processes multiple repositories from a YAML config file.

```bash
release-foundry -config repos.yml -days 7 -render github-release -out ./out
```

Output: a single `release-summary.json` with all repos, and one rendered file per
renderer in the output directory.

### Topic discovery mode (`-topic`)

Auto-discovers repositories in an org by GitHub topic, then processes them as a batch.

```bash
release-foundry -topic active -owner myorg -render github-release -out ./out

# With a config file to supply owner/baseBranch defaults
release-foundry -topic active -config repos.yml -render github-release -out ./out
```

The `-owner` flag (or `GITHUB_OWNER` env) specifies which org to search. If `-config`
is also provided, the config's `defaults.owner` fills in if `-owner` is not set.

## Batch config format (repos.yml)

```yaml
defaults:
  owner: myorg          # applied to all repos that don't set their own owner
  baseBranch: main      # default base branch

repos:
  - repo: api-server

  - repo: platform
    edition: enterprise   # included in JSON output; useful for segmenting

  - repo: plugin-sdk
    baseBranch: release/v2
    includeLabels: [feature, fix, plugin]
    excludeLabels: [internal]
```

## Filtering rules

A PR is **included** if:

1. It is merged (not open or draft)
2. Its base branch matches `baseBranch` (default: `main`)
3. Its merge date falls within the time window (`-since` or `-days`)
4. It has at least one of the include labels **or** its title matches a conventional
   commit prefix that maps to an included type
5. It does not have an exclude label **and** its title does not match an excluded prefix

Default include labels: `feature`, `fix`, `performance`, `security`, `infrastructure`

Default exclude labels: `internal`, `refactor`, `chore`

Excluded conventional commit prefixes: `chore`, `refactor`, `docs`, `ci`, `test`,
`style`

## Output schema

### Single-repo JSON (`release-summary.json`)

```json
{
  "generated_at": "2024-04-15T12:00:00Z",
  "repository": "myorg/myrepo",
  "edition": "enterprise",
  "time_window_days": 7,
  "summary_stats": {
    "total_prs": 12,
    "features": 4,
    "fixes": 6,
    "performance": 2
  },
  "pull_requests": [
    {
      "number": 42,
      "type": "feature",
      "title": "feat: add multi-repo batch mode",
      "customer_impact_raw": "Teams can now process all repos in one command.",
      "technical_summary": "Implemented batch YAML config loader and parallel collection.",
      "metrics": "Reduced release prep time from 30 min to 2 min.",
      "marketing_notes": "First tool to support org-wide release notes in one pass.",
      "labels": ["feature"],
      "author": "alice",
      "files_changed": 8,
      "additions": 320,
      "deletions": 45,
      "merged_at": "2024-04-14T09:30:00Z"
    }
  ]
}
```

### Batch JSON

Wraps an array of single-repo summaries:

```json
{
  "generated_at": "2024-04-15T12:00:00Z",
  "time_window_days": 7,
  "repositories": [
    { /* WeeklySummary for repo 1 */ },
    { /* WeeklySummary for repo 2 */ }
  ]
}
```

## Available renderers

| Name | Output file | Description |
|------|-------------|-------------|
| `github-release` | `<repo>-github-release.md` | Markdown grouped by PR type with emoji headers. Suitable for GitHub release body. |

Additional renderers can be added. See [CONTRIBUTING.md](../CONTRIBUTING.md) for the
renderer interface.

## Environment variables

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | GitHub personal access token |
| `GITHUB_OWNER` | Default repository owner or org |
| `GITHUB_REPO` | Default repository name |

## Examples

```bash
# Last 30 days, single repo
release-foundry -owner golang -repo go -days 30 -render github-release -out ./out

# Since a specific tag date
release-foundry -owner myorg -repo myrepo -since 2024-01-15T00:00:00Z \
  -render github-release -out ./out

# Batch: all repos in repos.yml, last 2 weeks
release-foundry -config repos.yml -days 14 -render github-release -out ./out

# Topic discovery: all repos tagged "release-managed" in the org
release-foundry -topic release-managed -owner myorg -render github-release -out ./out

# JSON only (no renderer), custom output path
release-foundry -owner myorg -repo myrepo -output ./data/summary.json

# Print version
release-foundry version
```
