# release-foundry

Collects merged pull request data from a GitHub repository and outputs a structured JSON summary for downstream LLM-based marketing content generation.

## Prerequisites

- Go 1.22+
- A GitHub personal access token with `repo` scope (or fine-grained token with pull-request read access)

## Setup

```bash
make build
```

## Configuration

Each of `token`, `owner`, and `repo` is resolved in order:

1. **CLI flag** (`-token`, `-owner`, `-repo`)
2. **Environment variable** (`GITHUB_TOKEN`, `GITHUB_OWNER`, `GITHUB_REPO`)
3. **Interactive prompt** (if neither flag nor env var is set; token input is masked)

### Option A — CLI flags

```bash
./release-foundry -token ghp_... -owner your-org -repo your-repo
```

### Option B — Environment variables

```bash
export GITHUB_TOKEN="ghp_..."
export GITHUB_OWNER="your-org"
export GITHUB_REPO="your-repo"
./release-foundry
```

### Option C — Interactive prompts

```bash
./release-foundry
# GitHub token: ••••••••
# GitHub owner: your-org
# GitHub repo:  your-repo
```

You can mix and match — e.g. set the token via env var and pass owner/repo as flags.

## Options

| Flag      | Default                              | Description                                      |
|-----------|--------------------------------------|--------------------------------------------------|
| `-token`  | `$GITHUB_TOKEN`                      | GitHub personal access token                     |
| `-owner`  | `$GITHUB_OWNER`                      | GitHub repository owner                          |
| `-repo`   | `$GITHUB_REPO`                       | GitHub repository name                           |
| `-days`   | `7`                                  | Number of days to look back                      |
| `-output` | `weekly_engineering_summary.json`    | Output file path                                 |
| `-config` | *(none)*                             | Path to multi-repo YAML config for batch mode    |

Example — last 14 days, custom output path:

```bash
./release-foundry -owner my-org -repo my-repo -days 14 -output summary.json
```

## Batch Mode (Multi-Repo)

Use `-config` to process multiple repositories in a single run. When `-config` is set, the `-owner` and `-repo` flags are ignored — repos are defined in the YAML file.

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
- `edition` tags each repo's output (e.g. `"enterprise"`, `"free"`) so downstream consumers know which edition a PR belongs to.
- `includeLabels` / `excludeLabels` override the default filter rules per repo.

See `repos.example.yml` for a working example.

### Batch output

```json
{
  "generated_at": "2025-01-15T10:00:00Z",
  "time_window_days": 7,
  "repositories": [
    {
      "repository": "your-org/your-repo",
      "edition": "enterprise",
      "time_window_days": 7,
      "summary_stats": { "total_prs": 5, "features": 3, "fixes": 2, "performance": 0 },
      "pull_requests": [...]
    },
    {
      "repository": "your-org/your-repo-free",
      "edition": "free",
      ...
    }
  ]
}
```

## Filtering

PRs must satisfy all of:
1. Merged into `main`
2. Merged within the configured time window
3. Has at least one include label: `feature`, `fix`, `performance`, `security`, `infrastructure`
4. Does **not** have any exclude label: `internal`, `refactor`, `chore`

## Output

The tool writes `weekly_engineering_summary.json` (or the path given via `-output`) with this structure:

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

If the PR body contains a `## Customer Impact` section, its content is extracted into `customer_impact_raw`. Otherwise the field is empty. No AI summarization is performed — this service only collates structured data.

## Project Structure

```
cmd/release-foundry/main.go    # Entrypoint, config loading, CLI flags
internal/
  config/config.go             # YAML loader for multi-repo batch config
  domain/model.go              # Domain types, label rules, classification
  domain/filter.go             # Configurable filter rules (labels, title prefixes)
  github/client.go             # GitHub REST API client with pagination & rate-limit handling
  service/collector.go         # Orchestration: fetch, filter, enrich, structure
repos.example.yml              # Example multi-repo batch config
```
