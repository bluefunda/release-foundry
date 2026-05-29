# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |
| Older   | No        |

We support only the latest released version. Please upgrade before reporting.

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

To report a vulnerability, please open a
[private security advisory](../../security/advisories/new) on GitHub.
Include:

- A description of the vulnerability and its potential impact
- Steps to reproduce the issue
- Any proof-of-concept code or screenshots

We will respond within **72 hours** and aim to ship a patch within **14 days**
for confirmed vulnerabilities.

## Scope

In scope:

- Token handling and credential exposure (`internal/github/`, `cmd/release-foundry/`)
- GitHub API request/response handling
- YAML config parsing (path traversal, injection)
- CLI argument handling

Out of scope:

- Issues in upstream dependencies (report to the upstream project)
- GitHub Actions runner security (GitHub's responsibility)
- Rate limiting as a denial-of-service vector (no user data at risk)

## Security Considerations for Users

- **Token scopes**: release-foundry requires read access to the repositories it collects
  data from. Use a fine-grained PAT scoped to `contents: read` and `pull-requests: read`
  where possible.
- **Token storage**: never pass tokens via command-line arguments in shared environments
  (they appear in process listings). Use the `GITHUB_TOKEN` env var or the `gh` CLI.
- **Config files**: `repos.yml` may contain org/repo names. Do not commit tokens or
  secrets into config files.
- **Output files**: generated JSON output (`release-summary.json`) contains PR titles
  and authors but no secrets. Treat it as public if your repo is public.
