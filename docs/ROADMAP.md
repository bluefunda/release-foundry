# release-foundry roadmap

## Vision

release-foundry is a **release content pipeline** — a tool that sits between code
changes and external communication. It transforms developer-authored PR data
(the `Customer Impact`, `Metrics`, and `Marketing Notes` fields from PR templates)
into audience-specific artifacts: changelogs, release notes, blog posts, social
posts, newsletters, and internal digests.

The current implementation (PR harvester → structured JSON → pluggable renderers)
is the foundation. The roadmap extends it into a **multi-renderer pipeline** with
channel-specific outputs and an optional generative layer.

## Operating model

### 1. Technical team operates on changes

The PR template + conventional commits are the input contract:

- **PR title format enforced by CI** (`feat:`, `fix:`, `perf:`, `security:`, `chore:`)
  — release-please needs this to auto-version
- **Labels auto-applied** from PR title via action — drives release-foundry filtering
- **release-please runs on every push to main** — produces a release PR with bumped
  version + per-repo `CHANGELOG.md`; merging the release PR creates the tag
- **release-foundry runs on tag creation** — picks up the new tag and produces all
  downstream artifacts in one pass

The team's only manual step: write a good PR body. Everything else is automatic.

### 2. Multi-renderer architecture

One PR-data ingest, many output formats:

| Renderer | Output | Channel |
|---|---|---|
| `github-release` | Markdown grouped by type with emoji headers | GitHub release body |
| `changelog` | Standard CHANGELOG format | Repo changelog file |
| `blog-post` | LLM-drafted long-form using `Customer Impact` + `Metrics` | Blog repo PR |
| `linkedin-post` | 3-5 paragraph customer-impact story | Buffer / draft file |
| `tweet-thread` | 5-tweet thread, 1 per top feature | X API draft |
| `newsletter` | Monthly digest aggregating all repos | Email platform |
| `internal-digest` | Slack/email "what shipped this week" | Webhook |
| `rss` | Atom/RSS feed aggregating releases | Feed reader |

Two distinct layers:

- **Deterministic templates** (changelog, GitHub release body) — no LLM, fact-safe,
  never misrepresents
- **LLM-generated narratives** (blog, social) — drafts only, human review required
  before publish

CLI surface (planned):

```bash
release-foundry run \
  --repos repos.yml \
  --since v1.2.0 \
  --render github-release,linkedin-post,blog-post \
  --out ./release-artifacts/ \
  --draft   # never auto-publish; --publish is opt-in
```

### 3. Community engagement extensions

- **Cross-repo intelligence** — multi-repo batch already works; add pattern
  detection ("two repos both shipped SSO this month" → bundled announcement)
- **Contributor recognition** — list external contributors with PR links in every
  release artifact
- **Customer impact rollup** — quarterly artifact aggregating every PR's
  `Customer Impact` field into a "what we shipped for you" doc
- **Roadmap surface** — pull open issues with `roadmap` label across repos, render
  to a public roadmap page

## Phasing

| Phase | Scope | Status |
|---|---|---|
| **1** | PR harvest → `github-release` renderer → GitHub release body automation | ✅ Done |
| **2** | Topic discovery, batch mode, renderer registry | ✅ Done |
| **3** | `blog-post` and `linkedin-post` renderers (draft only) + RSS feed | Planned |
| **4** | Auto-publishing pipeline + webhook outputs + roadmap surface | Planned |
| **5** | Plugin system, contributor recognition, customer rollups | Future |

## Future architecture

- **Plugin system** — load external renderers from a shared library or binary plugin
- **Templating** — user-defined Mustache/Go templates for custom output formats
- **Web UI** — browser-based dashboard for browsing releases and triggering renders
- **Git provider abstraction** — support GitLab, Gitea alongside GitHub
- **Release analytics** — track velocity, contributor diversity, PR cycle times
- **AI-assisted changelog** — LLM summarization layer with attribution and citation

## What this is NOT

- Not a metrics dashboard running continuously — it's a CLI run on demand or on tag
- Not a CI tool — it consumes CI output (tags, PR data); it doesn't run tests or
  build code
- Not a CMS — drafts go to files or existing platforms; release-foundry doesn't host
  content

## How to act on this roadmap

When ready to build a phase:

1. Extract that phase into a concrete spec (file paths, function signatures, acceptance
   criteria) in a GitHub issue
2. Implement
3. Update this ROADMAP.md to mark the phase complete and refine subsequent phases
