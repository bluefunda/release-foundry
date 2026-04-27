# release-foundry roadmap

## Vision

release-foundry is the **content hub of bluefunda** — the pipeline that sits between code changes and external communication. It transforms developer-authored PR data (the existing PR template's `Customer Impact`, `Metrics`, `Marketing Notes` fields) into audience-specific artifacts: changelogs, release notes, blog posts, social posts, newsletters, and internal digests.

The current implementation (PR harvester → structured JSON → LLM-friendly input) is the foundation. The roadmap extends it into a **multi-renderer pipeline** with channel-specific outputs.

## Operating model

### 1. Technical team operates on changes

The PR template + conventional commits are the input contract. Tighten the loop:

- **PR title format enforced by CI** (`feat:`, `fix:`, `perf:`, `security:`, `chore:`) — release-please needs this to auto-version
- **Labels auto-applied** from PR title via a small action — drives release-foundry's filtering
- **release-please runs on every push to main** in each repo — produces a release PR with bumped version + per-repo `CHANGELOG.md`. Merging the release PR creates the tag.
- **release-foundry runs on tag creation** — picks up the new tag and produces all downstream artifacts in one go.

The team's only manual step: write a good PR body. Everything else is automatic.

### 2. Structured content generation — multi-renderer architecture

Add a `renderers/` directory with templates per channel. One PR-data ingest, many output formats:

| Renderer | Output | Channel | Trigger |
|---|---|---|---|
| `github-release` | Markdown grouped by type with emoji headers | `gh release create` per repo | Per tag |
| `blog-post` | LLM-drafted long-form using `Customer Impact` + `Metrics` | PR to `bluefunda.com` repo | Major releases |
| `linkedin-post` | 3-5 paragraph customer-impact story | Buffer/Hypefury API or draft file | Per release |
| `tweet-thread` | 5-tweet thread, 1 per top feature, with screenshots | X API draft | Per release |
| `devto-article` | Technical deep-dive, code samples included | Dev.to API draft | Major features only |
| `newsletter` | Monthly digest aggregating all repos | Buttondown/Substack | Cron, monthly |
| `internal-digest` | Slack/email "what shipped this week" | Webhook | Cron, weekly |

Two distinct layers:

- **Deterministic templates** (changelog, GitHub release body) → no LLM, fact-safe, never misrepresents
- **LLM-generated narratives** (blog, social) → drafts only, human-review required before publish

CLI surface:

```bash
release-foundry run \
  --repos repos.yml \
  --since v1.2.0 \
  --render github-release,linkedin-post,blog-post \
  --out ./release-artifacts/ \
  --draft   # never auto-publish; --publish flips this
```

### 3. Community engagement extensions

Beyond release notes, release-foundry can drive sustained community engagement:

**Cross-repo intelligence** — multi-repo batch is already there. Use it to spot patterns:
- "ABAPer + cai-cli both shipped SSO this month" → bundled announcement
- "btp-go's connectivity API powered three customer integrations" → case study draft

**Roadmap surface** — pull open issues with `roadmap` label across repos, render to a public roadmap page in the website. Updated on every release.

**Contributor recognition** — every release artifact lists external contributors with PR links. Drives community participation; cheap and high-signal.

**Customer impact rollup** — quarterly artifact aggregating every PR's `Customer Impact` field into a "what we shipped for you" doc, segmented by customer label (`customer:acme`, `customer:globex`).

**Partner/integration changelog** — filter to PRs touching public APIs or integrations, send to partners on a different cadence than the public changelog.

**RSS/Atom feed** — aggregate GitHub releases across all bluefunda repos into one feed at `bluefunda.com/feed.xml`. Free SEO + lets developers subscribe in their reader of choice.

**Webhook outputs** — Discord/Slack hooks for "X just shipped" — community channel engagement, free.

## Phasing

| Phase | Scope | Outcome |
|---|---|---|
| **1 (this week)** | PR title CI check + release-please in every repo + `github-release` renderer | Every tag creates a real GitHub release automatically |
| **2 (this month)** | `blog-post` and `linkedin-post` renderers (draft only) + internal weekly digest | Human still publishes, but drafts arrive automatically |
| **3 (next quarter)** | Auto-publishing pipeline + RSS feed + roadmap page | Releases reach external audiences without human in the loop |
| **4 (later)** | Customer rollups + partner channel + contributor recognition | Sustained community engagement loop |

The critical move is **Phase 1**. Once tags reliably produce GitHub releases via release-foundry, every subsequent renderer is a small additive change.

## What this is NOT

- Not a metrics dashboard or analytics platform — keep it as a CLI run on demand or on tag, not a service running continuously
- Not a CI tool — it consumes CI output (tags, PR data); it doesn't run tests or build code
- Not a CMS — drafts are written to files or pushed to existing platforms (GitHub, Buttondown, etc.); release-foundry doesn't host content

## How to act on this roadmap

When ready to build a phase:

1. Extract that phase into a `PROMPTS.md` build spec (concrete file paths, function signatures, acceptance criteria)
2. Implement
3. Update this ROADMAP.md to mark the phase complete and refine subsequent phases based on what was learned
