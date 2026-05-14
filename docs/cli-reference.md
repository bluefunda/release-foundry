# CLI Reference

## Prerequisites

- Go 1.22+
- A GitHub personal access token with `repo` scope (or fine-grained token with pull-request read access)

## Build

```bash
make build
```

## Configuration

Each of `token`, `owner`, and `repo` is resolved in order:

1. **CLI flag** (`-token`, `-owner`, `-repo`)
2. **Environment variable** (`GITHUB_TOKEN`, `GITHUB_OWNER`, `GITHUB_REPO`)
3. **Interactive prompt** (token input is masked)

```bash
# Flags
./release-foundry -token ghp_... -owner your-org -repo your-repo

# Environment variables
export GITHUB_TOKEN="ghp_..."
export GITHUB_OWNER="your-org"
export GITHUB_REPO="your-repo"
./release-foundry

# Mix — token from env, owner/repo as flags
./release-foundry -owner your-org -repo your-repo
```

## Flags

| Flag      | Default                           | Description                                   |
|-----------|-----------------------------------|-----------------------------------------------|
| `-token`  | `$GITHUB_TOKEN`                   | GitHub personal access token                  |
| `-owner`  | `$GITHUB_OWNER`                   | GitHub repository owner                       |
| `-repo`   | `$GITHUB_REPO`                    | GitHub repository name                        |
| `-days`   | `7`                               | Number of days to look back                   |
| `-since`  | *(none)*                          | RFC3339 timestamp — overrides `-days`         |
| `-output` | `weekly_engineering_summary.json` | Output file path                              |
| `-render` | *(none)*                          | Renderer to use (e.g. `github-release`)       |
| `-out`    | *(none)*                          | Output directory for rendered files           |
| `-config` | *(none)*                          | Path to multi-repo YAML config (batch mode)   |

```bash
# Last 14 days, custom output
./release-foundry -owner my-org -repo my-repo -days 14 -output summary.json

# Since a specific timestamp
./release-foundry -owner my-org -repo my-repo -since 2026-04-01T00:00:00Z -render github-release -out ./out
```

## Batch Mode

Use `-config` to process multiple repositories in one run. `-owner` and `-repo` flags are ignored when `-config` is set.

```bash
./release-foundry -config repos.yml -days 7 -output batch.json
```

### Config file format

```yaml
defaults:
  owner: your-org
  baseBranch: main

repos:
  - repo: your-repo
    edition: enterprise

  - repo: your-repo-free
    edition: free

  - repo: your-plugins
    edition: enterprise
    includeLabels: [feature, fix, plugin]
    excludeLabels: [internal, wip]
```

- `defaults.owner` and `defaults.baseBranch` apply to all repos unless overridden per entry.
- `edition` tags each repo's output (e.g. `"enterprise"`, `"free"`).
- `includeLabels` / `excludeLabels` override the default filter rules per repo.

See `repos.example.yml` for a working example.

## Filtering

PRs must satisfy all of:

1. Merged into `main`
2. Merged within the configured time window
3. Has at least one include label: `feature`, `fix`, `performance`, `security`, `infrastructure`
4. Does **not** have any exclude label: `internal`, `refactor`, `chore`

## Output Format

```json
{
  "generated_at": "2025-01-15T10:00:00Z",
  "repository": "your-org/your-repo",
  "time_window_days": 7,
  "summary_stats": {
    "total_prs": 12,
    "features": 5,
    "fixes": 4,
    "performance": 3
  },
  "pull_requests": [
    {
      "number": 456,
      "type": "feature",
      "title": "Add SSO login",
      "customer_impact_raw": "Users can now sign in with SSO.",
      "technical_summary": "Full PR body text...",
      "labels": ["feature"],
      "author": "engineer",
      "files_changed": 12,
      "additions": 340,
      "deletions": 20,
      "merged_at": "2025-01-14T18:30:00Z"
    }
  ]
}
```

If the PR body contains a `## Customer Impact` section, its content is extracted into `customer_impact_raw`. Otherwise the field is empty. No AI summarization is performed — this tool only collates structured data.

### Batch output

```json
{
  "generated_at": "2025-01-15T10:00:00Z",
  "time_window_days": 7,
  "repositories": [
    {
      "repository": "your-org/your-repo",
      "edition": "enterprise",
      "summary_stats": { "total_prs": 5, "features": 3, "fixes": 2 },
      "pull_requests": [...]
    }
  ]
}
```
