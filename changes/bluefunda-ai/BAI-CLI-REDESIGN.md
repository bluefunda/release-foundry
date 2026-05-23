# `bai` CLI — Architecture Redesign

> **Status:** Approved for implementation
> **Target repo:** `bluefunda-ai`
> **Implementation branch:** `feat/cli-redesign`

---

## Executive Summary

`bai` today is an enterprise platform SDK that happens to have a chat feature. The goal is to invert that: make `bai` an AI pair programmer that happens to have platform plumbing underneath. The structural change is decisive but surgical — the gRPC machinery, auth, TUI, and streaming pipeline stay almost entirely intact. What changes is the surface area, the command routing, and the conceptual model presented to users.

The single most important change: **running `bai` with no arguments should open an interactive AI session** — not print help text.

---

## CLI Philosophy

| Old | New |
|---|---|
| "CLI for the BlueFunda AI platform" | "Your AI pair programmer" |
| User must know about BFF, gRPC, realms | Invisible infrastructure |
| `bai chat start` to talk to AI | `bai` to talk to AI |
| 10 top-level commands | 6 visible, rest hidden |
| Enterprise SDK terminology | Developer-first language |
| Flags like `--bff`, `--gateway` | Auto-configured, hidden |
| Help output: exhaustive | Help output: minimal, inviting |

---

## New Command Tree

```
bai                         ← interactive AI session (REPL)
bai [prompt]                ← start session with initial message
bai code [prompt]           ← agentic coding mode
bai login                   ← sign in
bai doctor                  ← diagnose environment
bai config                  ← configuration
  bai config get [key]
  bai config set key=value
  bai config list
bai mcp                     ← tool integrations
  bai mcp list
  bai mcp add <name>
  bai mcp remove <name>
bai version

── Hidden (not in --help, but functional) ──────────────────────────
bai sessions                ← list/manage past sessions
bai update                  ← self-update binary
bai _debug                  ← raw gRPC/auth diagnostics
bai _admin                  ← account/billing/rate-limit

── Removed from top-level ──────────────────────────────────────────
bai chat         → bai (implicit)
bai health       → bai doctor
bai model        → /model slash command inside session
bai user         → /account slash command inside session
bai billing      → /usage slash command inside session
bai rate-limit   → /usage slash command inside session
```

### Cobra tree (Go)

```go
rootCmd         // bai [prompt] — REPL or arg passthrough
├── codeCmd     // bai code [prompt]
├── loginCmd    // bai login
├── doctorCmd   // bai doctor
├── configCmd   // bai config
│   ├── configGetCmd
│   ├── configSetCmd
│   └── configListCmd
├── mcpCmd      // bai mcp
│   ├── mcpListCmd
│   ├── mcpAddCmd
│   └── mcpRemoveCmd
├── versionCmd  // bai version
│
│── Hidden (cobra.Command.Hidden = true) ───────────────────────────
├── sessionsCmd     // bai sessions [list|resume|delete]
├── updateCmd       // bai update
├── debugCmd        // bai _debug
└── adminCmd        // bai _admin [user|billing|rate-limit|health]
    ├── adminUserCmd
    ├── adminBillingCmd
    ├── adminRateLimitCmd
    └── adminHealthCmd
```

---

## Command Migration Table

| Current Command | Disposition | New Location |
|---|---|---|
| `bai chat start` | **Merge** → root | `bai` / `bai "prompt"` |
| `bai chat start --new` | **Merge** | `bai --new` |
| `bai chat list` | **Hide** | `bai sessions` (hidden) |
| `bai chat history <id>` | **Hide** | `/history` slash cmd |
| `bai chat context <id>` | **Hide** | `bai _debug context` |
| `bai chat title <id>` | **Remove** | auto-generated, not user-facing |
| `bai chat stop <id>` | **Remove** | Ctrl+C in session |
| `bai code` | **Keep**, promote | `bai code` (flagship) |
| `bai health` | **Rename** | `bai doctor` |
| `bai login` | **Keep** | `bai login` |
| `bai version` | **Keep** | `bai version` |
| `bai model list` | **Hide** | `/model` slash cmd |
| `bai mcp list` | **Keep** | `bai mcp list` |
| `bai mcp user` | **Merge** | `bai mcp list` (combined view) |
| `bai mcp select` | **Rename** | `bai mcp add <name>` |
| `bai user info` | **Hide** | `/account` slash cmd |
| `bai user settings` | **Hide** | `bai _admin user` |
| `bai billing subscription` | **Hide** | `/usage` slash cmd |
| `bai billing plans` | **Hide** | `bai _admin billing` |
| `bai rate-limit` | **Hide** | `/usage` slash cmd |
| `--bff` flag | **Hide** | `bai config set endpoint=...` |
| `--gateway` flag | **Hide** | `bai config set gateway=...` |
| `--domain` flag | **Hide** | `bai config set domain=...` |

---

## Interactive UX Design

### Default session startup (`bai` or `bai "prompt"`)

```
$ bai

  BlueFunda AI  ·  claude-sonnet  ·  📁 bluefunda-ai
  ─────────────────────────────────────────────────────
  Resuming session a3f7c821  ·  "fix failing tests"

  You
  explain the auth middleware

  Assistant
  The auth middleware in internal/grpc/conn.go handles...

  ─────────────────────────────────────────────────────
  > █
  Enter send · Shift+Enter newline · / commands · Ctrl+C quit
```

Key behaviors on startup:

1. **Detect git repo** — show repo name in header
2. **Detect language/framework** — inject as context hint
3. **Resume last session** — if last session exists and is recent (< 24h), offer resume
4. **Show model** — always visible in header
5. **No boilerplate** — no walls of welcome text

### Agentic coding session (`bai code`)

```
$ bai code "add input validation to the user registration endpoint"

  BlueFunda AI  ·  code  ·  claude-opus  ·  📁 /src/bluefunda-ai
  ─────────────────────────────────────────────────────────────────
  You
  add input validation to the user registration endpoint

  Assistant  ⠙
  ● read_file  src/handlers/user.go
  ✓  read_file  0.1s  —  312 lines

  ● write_file  src/handlers/user.go
  Apply? [y/N] █
```

---

## Slash Command System

Slash commands work inside any interactive session (`bai` and `bai code`).

```
Session control
  /new              Start a new session
  /resume [id]      Resume a past session (shows list if no id)
  /sessions         List recent sessions with titles
  /clear            Clear conversation display
  /reset            New session in current window
  /exit, /quit      Exit

AI control
  /model [name]     Show current model or switch model
  /tools            List available tools
  /mcp [name]       Show MCP servers or activate one
  /context          Show session context (chat ID, message count)

Workflow
  /code             Switch to code mode (load file system tools)
  /chat             Switch to chat mode (drop file tools)
  /auto             Toggle auto-apply for code mode

Account (lazy-loaded)
  /account          Show name, email, plan
  /usage            Show token usage and rate limits
  /plans            Show available plans

Meta
  /help             Show this list
  /debug            Show session debug info (hidden)
```

---

## Root Command Behavior

### Today

```go
var rootCmd = &cobra.Command{
    Use:   "bai",
    Short: "bai -- CLI for the BlueFunda AI platform",
    // No RunE — prints help
}
```

### After redesign

```go
var rootCmd = &cobra.Command{
    Use:   "bai [prompt]",
    Short: "Your AI pair programmer",
    Args:  cobra.ArbitraryArgs,
    RunE:  runDefault,
}

func runDefault(cmd *cobra.Command, args []string) error {
    var initialPrompt string
    if len(args) > 0 {
        initialPrompt = strings.Join(args, " ")
    }
    return runChatSession(initialPrompt, forceNew)
}
```

Cobra's subcommand routing means named subcommands (`code`, `login`, etc.) still resolve correctly. Unrecognised args reach `runDefault` and become the initial prompt.

---

## `bai doctor` (replaces `bai health`)

```
$ bai doctor

  Checking your environment...

  ✓  Config file          ~/.bai/config.yaml
  ✓  Authentication       token valid (expires in 47m)
  ✓  Backend reachable    cli.bluefunda.com:443  (28ms)
  ✓  Account              phani.p@bluefunda.com
  ✗  Model default        not set — using "openai" fallback
  ℹ  MCP servers          none selected

  1 warning · run `bai config set model=claude-sonnet` to fix
```

Check sequence:

1. Config file present and parseable
2. Token present and not expired
3. Backend gRPC ping
4. `GetUserInfo` call (validates auth end-to-end)
5. Default model configured
6. MCP server selected (informational)

---

## `bai config`

```bash
bai config list
bai config get model
bai config set model=claude-sonnet
bai config set endpoint=cli.bluefunda.com:443
```

### Config key mapping

| User-facing key | Internal field | Default |
|---|---|---|
| `endpoint` | `BFFURL` | `cli.bluefunda.com:443` |
| `model` | `Defaults.Model` | `openai` |
| `output` | `Defaults.Output` | `table` |
| `gateway` | `GatewayURL` | hidden/advanced |
| `domain` | `Domain` | hidden/advanced |
| `realm` | `Realm` | hidden/advanced |

### Config file evolution

```yaml
# Before
gatewayurl: https://ai.bluefunda.com
bffurl: cli.bluefunda.com:443
domain: bluefunda.com
realm: trm
auth:
  accesstoken: eyJ...
  refreshtoken: eyJ...
  tokenexpiry: 2026-05-23T...
defaults:
  model: openai
  output: table

# After
endpoint: cli.bluefunda.com:443
model: claude-sonnet
auth:
  access_token: eyJ...
  refresh_token: eyJ...
  token_expiry: 2026-05-23T...
# advanced fields (omitted unless set)
# gateway: https://ai.bluefunda.com
# domain: bluefunda.com
# realm: trm
```

Migration runs once inside `config.Load()` — reads old field names, writes new ones.

---

## Help Output

### Before

```
bai is a command-line interface for interacting with the BlueFunda AI platform via gRPC.

Usage:
  bai [command]

Available Commands:
  billing     Billing and subscription operations
  chat        Chat operations
  code        Agentic coding session with local file system access
  completion  Generate the autocompletion script for the specified shell
  health      Check BFF connectivity (gRPC)
  help        Help about any command
  login       Log in via browser (opens Keycloak)
  mcp         MCP server management
  model       LLM model operations
  rate-limit  Query current rate limit status
  user        User account operations
  version     Print the CLI version

Flags:
      --bff string       BFF gRPC address host:port (overrides config)
      --domain string    Domain (overrides config)
      --gateway string   Gateway base URL (overrides config)
  -h, --help             help for bai
  -o, --output string    Output format: table, json, quiet
      --version          version for bai
```

### After

```
BlueFunda AI — your AI pair programmer

Usage:
  bai [prompt]
  bai [command]

Examples:
  bai                          start interactive session
  bai "fix the failing tests"  start with a message
  bai code                     agentic coding mode
  bai login                    sign in

Commands:
  code      Agentic coding workflow
  login     Sign in
  doctor    Diagnose your environment
  config    Configure CLI settings
  mcp       Manage tool integrations
  version

Use /help inside a session to see available slash commands.
```

---

## Example User Sessions

### First-time user

```bash
$ bai
  Not signed in. Run `bai login` to get started.

$ bai login
  Opening browser for sign-in...
  ✓ Signed in as phani.p@bluefunda.com

$ bai
  BlueFunda AI  ·  claude-sonnet  ·  📁 bluefunda-ai
  > █
```

### Developer daily workflow

```bash
$ bai "why are the integration tests failing?"

$ bai code "fix the failing auth tests and open a PR"

$ bai code
  > read src/handlers/user.go and add input validation
  > /model claude-opus
  > run the tests and fix any failures
```

### Power user / scripts (hidden commands still work)

```bash
bai _admin user
bai _admin billing subscription -o json
bai _admin rate-limit -o json | jq '.remaining'
bai --bff localhost:50051 code "debug this locally"
```

---

## Go Package Structure

```
cmd/
  bai/
    main.go

internal/
  cmd/
    root.go          ← bai [prompt] — default REPL handler  (was: no RunE)
    code.go          ← bai code                             (unchanged)
    login.go         ← bai login                            (unchanged)
    doctor.go        ← bai doctor                           (new; replaces health.go)
    config.go        ← bai config get/set/list              (new)
    mcp.go           ← bai mcp list/add/remove              (simplified)
    version.go       ← bai version                          (unchanged)
    sessions.go      ← bai sessions (hidden)                (new; was chat list)
    admin.go         ← bai _admin (hidden)                  (new; wraps user/billing/rl)
    helpers.go       ← bffConn, printer, etc.               (unchanged)
    chat.go          ← kept hidden for backward compat       (hidden)
    health.go        ← kept hidden, delegates to doctor      (hidden)

  session/           ← new package; session lifecycle
    manager.go       ← create, resume, list sessions

  agent/             ← extracted from code.go
    loop.go          ← 20-iteration tool dispatch loop
    tools.go         ← tool schemas
    filesystem.go    ← tool execution

  tui/               ← mostly unchanged
    model.go
    view.go
    stream.go
    slashcmd.go      ← extended with /account /usage /model /resume /sessions
    messages.go
    theme.go
    codeblock.go
    program.go

  config/            ← simplified struct + one-time migration
    config.go

  auth/              ← unchanged
  grpc/              ← unchanged
  ui/                ← unchanged
```

---

## Implementation Phases

### Phase 1 — Surface reshape (high impact, low risk)

- `root.go`: Add `RunE: runDefault` that calls `runChatSession`
- `root.go`: New short/long/examples, new help output
- Mark `billing`, `user`, `rate-limit`, `chat`, `health`, `model` as `Hidden: true`
- Hide `--bff`, `--gateway`, `--domain` from help (keep functional)
- Add `doctor.go` (config + ping + auth + user info checks)
- Add `-o` / `--output` flag documentation cleanup

### Phase 2 — `bai code` promotion

- `bai code "prompt"` auto-submits initial message (mirror chat fix)
- `--auto` alias for `--auto-apply`
- Header shows working directory and detected repo name
- Auto-detect working dir from git root

### Phase 3 — Slash command expansion

- Extend `slashcmd.go`: `/model`, `/account`, `/usage`, `/resume`, `/sessions`, `/mcp`, `/code`, `/chat`, `/new`
- `/account` and `/usage` fire gRPC calls lazily, render inline in TUI
- `/model [name]` updates `m.cfg.Model` live
- `/sessions` renders numbered list; `/resume 3` resumes by index

### Phase 4 — `bai config` and config simplification

- Add `config.go` with `get`/`set`/`list` subcommands
- Migrate config field names in `config.Load()` (one-time, backward compatible)
- `bai config set model=...` is the new recommended way to set defaults

### Phase 5 — `bai doctor` and startup context

- Full 6-step doctor check sequence
- Startup: detect git repo, detect framework, offer resume if session < 24h old
- First-run: if no token, print `→ run bai login to get started`

### Phase 6 — Polish

- `bai update` — downloads latest binary from GitHub releases
- Shell completion (`bai completion bash/zsh/fish`)
- Session telemetry (anonymous, opt-in) — session logs in `~/.bai/sessions/`
- `bai _admin` consolidated hidden surface

---

## Backward Compatibility Strategy

All old commands remain functional. They disappear from `--help` only.

```go
// chat.go — kept, hidden
var chatCmd = &cobra.Command{Hidden: true, ...}

// health.go — delegates to doctor
var healthCmd = &cobra.Command{
    Hidden: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runDoctor(cmd, args)
    },
}
```

Deprecation timeline:

- **Phase 1–5:** Hidden but fully functional
- **v2.0:** Hidden commands print deprecation notice, still work
- **v3.0:** Remove `bai chat`, `bai user`, `bai billing`, `bai rate-limit`, `bai health`

---

## Risks and Tradeoffs

| Risk | Likelihood | Mitigation |
|---|---|---|
| Scripts using `bai chat start` break | Low | Command stays, just hidden |
| `bai "prompt"` ambiguous with subcommand routing | Low | Cobra resolves subcommands first; bare args only reach root if no match |
| Config migration corrupts tokens | Low | Read-old-write-new; backup on first migration |
| Enterprise users expect verbose command structure | Medium | Old commands remain hidden but functional via `bai _admin` |
| `bai` and `bai code` UX overlap | Medium | `bai code` = file system tools active; `bai` = no file tools, switchable via `/code` |

---

## Final Recommendations

**Do first (high impact, low risk):**

1. Make `bai` (no args) open the REPL
2. Hide `billing`, `rate-limit`, `user`, `health`, `model`, `chat` from `--help`
3. Add `bai doctor`
4. Remove BFF/gRPC/gateway/realm from all user-visible strings

**Do next:**

5. Expand slash commands to cover everything now hidden
6. Startup context detection (git repo, framework, resume candidate)
7. `bai config get/set/list`

**Do later:**

8. `bai update`
9. Shell completion
10. Session telemetry
11. `bai _admin` consolidated admin surface

**Don't do:**

- Don't remove the TUI — it's excellent
- Don't change the gRPC/streaming architecture
- Don't merge `bai` and `bai code` into one command — code mode with file tools is meaningfully different
- Don't add `bai chat` as visible top-level — the root command IS chat
