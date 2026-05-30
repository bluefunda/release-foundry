package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/release-foundry/internal/config"
	"github.com/release-foundry/internal/domain"
	gh "github.com/release-foundry/internal/github"
	"github.com/release-foundry/internal/renderers"
	"github.com/release-foundry/internal/service"
)

// buildVersion is injected at link time via -ldflags "-X main.buildVersion=v1.2.3".
var buildVersion = "dev"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Handle top-level subcommands before flag parsing so that
	// "release-foundry version" works without any flags.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("release-foundry %s\n", buildVersion)
			os.Exit(0)
		case "help":
			printUsage()
			os.Exit(0)
		}
	}

	token := flag.String("token", "", "GitHub personal access token (overrides GITHUB_TOKEN env var)")
	owner := flag.String("owner", "", "GitHub repository owner/org (overrides GITHUB_OWNER env var)")
	repo := flag.String("repo", "", "GitHub repository name (overrides GITHUB_REPO env var)")
	days := flag.Int("days", 7, "number of days to look back for merged PRs")
	sinceStr := flag.String("since", "", "fetch PRs merged after this RFC3339 timestamp (overrides -days)")
	output := flag.String("output", "release-summary.json", "output JSON file path")
	configPath := flag.String("config", "", "path to multi-repo YAML config for batch mode")
	topic := flag.String("topic", "", "GitHub topic filter for auto-discovering repos in an org (e.g. active)")
	renderFlag := flag.String("render", "", fmt.Sprintf("comma-separated renderers to run (available: %s)", strings.Join(renderers.Names(), ", ")))
	outDir := flag.String("out", ".", "output directory for rendered artifacts")

	flag.Usage = printUsage
	flag.Parse()

	renderNames := parseRenderFlag(*renderFlag)

	// Topic mode: discover repos by GitHub topic, then process as a batch.
	if *topic != "" {
		runTopicBatch(*configPath, *topic, *owner, *token, *days, *sinceStr, *output, renderNames, *outDir)
		return
	}

	// Batch mode: process multiple repos from config file.
	if *configPath != "" {
		runBatch(*configPath, *token, *days, *sinceStr, *output, renderNames, *outDir)
		return
	}

	// Single-repo mode (backward compatible).
	cfg, err := loadConfig(*token, *owner, *repo, *days, *sinceStr)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	client := gh.NewClient(cfg.Token)
	collector := service.NewCollector(client, cfg)

	summary, err := collector.Collect(context.Background())
	if err != nil {
		log.Fatalf("collection failed: %v", err)
	}

	if err := writeJSON(*output, summary); err != nil {
		log.Fatalf("write output: %v", err)
	}
	log.Printf("wrote %s (%d PRs)", *output, summary.SummaryStats.TotalPRs)

	if err := runRenderers(renderNames, *outDir, summary, nil); err != nil {
		log.Fatalf("render: %v", err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `release-foundry — GitHub PR-based release notes generator

Usage:
  release-foundry [flags]              single-repo mode
  release-foundry -config repos.yml   batch mode (multiple repos)
  release-foundry -topic active        topic discovery mode (auto-discover repos by GitHub topic)
  release-foundry version              print version

Flags:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Available renderers: %s

Examples:
  # Single repo, last 7 days
  release-foundry -owner myorg -repo myrepo -render github-release -out ./out

  # Single repo since a specific date
  release-foundry -owner myorg -repo myrepo -since 2024-01-01T00:00:00Z -render github-release

  # Batch mode from config file
  release-foundry -config repos.yml -days 14 -render github-release -out ./out

  # Auto-discover repos tagged "active" in an org
  release-foundry -topic active -owner myorg -render github-release -out ./out

Environment variables:
  GITHUB_TOKEN   GitHub personal access token
  GITHUB_OWNER   Default repository owner
  GITHUB_REPO    Default repository name

See docs/cli-reference.md for full documentation.
`, strings.Join(renderers.Names(), ", "))
}

func runBatch(configPath, flagToken string, days int, sinceStr string, output string, renderNames []string, outDir string) {
	repoCfg, err := config.LoadReposConfig(configPath)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Resolve token once for all repos.
	token, err := resolveToken(flagToken)
	if err != nil {
		log.Fatalf("token error: %v", err)
	}

	client := gh.NewClient(token)
	since, err := parseSince(sinceStr, days)
	if err != nil {
		log.Fatalf("invalid -since: %v", err)
	}

	batch := domain.BatchSummary{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		TimeWindowDays: days,
	}

	for _, entry := range repoCfg.Repos {
		log.Printf("processing %s/%s (edition=%s)", entry.Owner, entry.Repo, entry.Edition)

		var filters domain.FilterConfig
		if len(entry.IncludeLabels) > 0 || len(entry.ExcludeLabels) > 0 {
			filters = domain.NewFilterConfig(entry.IncludeLabels, entry.ExcludeLabels)
		}

		cfg := domain.Config{
			Token:      token,
			Owner:      entry.Owner,
			Repo:       entry.Repo,
			BaseBranch: entry.BaseBranch,
			WindowDays: days,
			Since:      since,
			Edition:    entry.Edition,
			Filters:    filters,
		}

		collector := service.NewCollector(client, cfg)
		summary, err := collector.Collect(context.Background())
		if err != nil {
			log.Printf("error collecting %s/%s: %v (skipping)", entry.Owner, entry.Repo, err)
			continue
		}

		batch.Repositories = append(batch.Repositories, *summary)
	}

	if batch.Repositories == nil {
		batch.Repositories = []domain.WeeklySummary{}
	}

	if err := writeJSON(output, batch); err != nil {
		log.Fatalf("write output: %v", err)
	}

	total := 0
	for _, r := range batch.Repositories {
		total += r.SummaryStats.TotalPRs
	}
	log.Printf("wrote %s (%d repos, %d total PRs)", output, len(batch.Repositories), total)

	if err := runRenderers(renderNames, outDir, nil, &batch); err != nil {
		log.Fatalf("render: %v", err)
	}
}

// runTopicBatch discovers repos in org tagged with topic, then processes them as a batch.
// If configPath is non-empty, defaults (owner, baseBranch) are read from the config file.
func runTopicBatch(configPath, topic, flagOwner, flagToken string, days int, sinceStr, output string, renderNames []string, outDir string) {
	token, err := resolveToken(flagToken)
	if err != nil {
		log.Fatalf("token error: %v", err)
	}

	client := gh.NewClient(token)
	since, err := parseSince(sinceStr, days)
	if err != nil {
		log.Fatalf("invalid -since: %v", err)
	}

	// Resolve owner and baseBranch: flag takes precedence; fall back to config defaults.
	orgOwner := flagOwner
	baseBranch := "main"
	if configPath != "" {
		defaults, err := config.LoadDefaults(configPath)
		if err != nil {
			log.Fatalf("config error: %v", err)
		}
		if orgOwner == "" {
			orgOwner = defaults.Owner
		}
		baseBranch = defaults.BaseBranch
	}
	if orgOwner == "" {
		orgOwner, err = resolve("", "GITHUB_OWNER", "GitHub owner", false)
		if err != nil {
			log.Fatalf("owner error: %v", err)
		}
	}

	log.Printf("discovering repos in org %q with topic %q", orgOwner, topic)
	repoNames, err := client.SearchReposByTopic(context.Background(), orgOwner, topic)
	if err != nil {
		log.Fatalf("topic search: %v", err)
	}
	if len(repoNames) == 0 {
		log.Fatalf("no repos found in org %q with topic %q", orgOwner, topic)
	}
	log.Printf("found %d repos", len(repoNames))

	batch := domain.BatchSummary{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		TimeWindowDays: days,
	}

	for _, name := range repoNames {
		log.Printf("processing %s/%s", orgOwner, name)
		cfg := domain.Config{
			Token:      token,
			Owner:      orgOwner,
			Repo:       name,
			BaseBranch: baseBranch,
			WindowDays: days,
			Since:      since,
		}
		collector := service.NewCollector(client, cfg)
		summary, err := collector.Collect(context.Background())
		if err != nil {
			log.Printf("error collecting %s/%s: %v (skipping)", orgOwner, name, err)
			continue
		}
		batch.Repositories = append(batch.Repositories, *summary)
	}

	if batch.Repositories == nil {
		batch.Repositories = []domain.WeeklySummary{}
	}

	if err := writeJSON(output, batch); err != nil {
		log.Fatalf("write output: %v", err)
	}

	total := 0
	for _, r := range batch.Repositories {
		total += r.SummaryStats.TotalPRs
	}
	log.Printf("wrote %s (%d repos, %d total PRs)", output, len(batch.Repositories), total)

	if err := runRenderers(renderNames, outDir, nil, &batch); err != nil {
		log.Fatalf("render: %v", err)
	}
}

// resolve returns the first non-empty value from: flag, env var, interactive prompt.
func resolve(flagVal, envKey, promptLabel string, secret bool) (string, error) {
	if flagVal != "" {
		return flagVal, nil
	}
	if v := os.Getenv(envKey); v != "" {
		return v, nil
	}
	return prompt(promptLabel, secret)
}

// resolveToken resolves the GitHub token: flag → env → `gh auth token` → interactive prompt.
func resolveToken(flagVal string) (string, error) {
	if flagVal != "" {
		return flagVal, nil
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		return v, nil
	}
	if t := ghAuthToken(); t != "" {
		log.Println("using token from `gh auth token`")
		return t, nil
	}
	return prompt("GitHub token", true)
}

func prompt(label string, secret bool) (string, error) {
	fmt.Fprintf(os.Stderr, "%s: ", label)
	if secret {
		fd := int(os.Stdin.Fd())
		if term.IsTerminal(fd) {
			b, err := term.ReadPassword(fd)
			fmt.Fprintln(os.Stderr) // newline after masked input
			if err != nil {
				return "", fmt.Errorf("read %s: %w", label, err)
			}
			v := strings.TrimSpace(string(b))
			if v == "" {
				return "", fmt.Errorf("%s is required", label)
			}
			return v, nil
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("%s is required", label)
	}
	v := strings.TrimSpace(scanner.Text())
	if v == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	return v, nil
}

// ghAuthToken shells out to `gh auth token` to retrieve the token from the GitHub CLI.
func ghAuthToken() string {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func parseSince(sinceStr string, days int) (time.Time, error) {
	if sinceStr != "" {
		return time.Parse(time.RFC3339, sinceStr)
	}
	return time.Now().UTC().AddDate(0, 0, -days), nil
}

func loadConfig(flagToken, flagOwner, flagRepo string, days int, sinceStr string) (domain.Config, error) {
	// Token resolution: flag → env → gh CLI → interactive prompt.
	token, err := resolveToken(flagToken)
	if err != nil {
		return domain.Config{}, err
	}
	owner, err := resolve(flagOwner, "GITHUB_OWNER", "GitHub owner", false)
	if err != nil {
		return domain.Config{}, err
	}
	repo, err := resolve(flagRepo, "GITHUB_REPO", "GitHub repo", false)
	if err != nil {
		return domain.Config{}, err
	}

	since, err := parseSince(sinceStr, days)
	if err != nil {
		return domain.Config{}, fmt.Errorf("invalid -since value: %w", err)
	}

	return domain.Config{
		Token:      token,
		Owner:      owner,
		Repo:       repo,
		BaseBranch: "main",
		WindowDays: days,
		Since:      since,
	}, nil
}

func parseRenderFlag(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for r := range strings.SplitSeq(s, ",") {
		if r = strings.TrimSpace(r); r != "" {
			out = append(out, r)
		}
	}
	return out
}

// runRenderers runs each named renderer and writes output files to outDir.
// Exactly one of summary/batch must be non-nil.
func runRenderers(names []string, outDir string, summary *domain.WeeklySummary, batch *domain.BatchSummary) error {
	if len(names) == 0 {
		return nil
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	for _, name := range names {
		r, ok := renderers.Get(name)
		if !ok {
			log.Printf("unknown renderer %q — skipping (available: %s)", name, strings.Join(renderers.Names(), ", "))
			continue
		}

		var content string
		var filename string

		if batch != nil {
			content = r.Batch(*batch)
			filename = fmt.Sprintf("%s.%s", name, r.FileExtension())
		} else if summary != nil {
			content = r.Single(*summary)
			repoSlug := strings.ReplaceAll(summary.Repository, "/", "-")
			filename = fmt.Sprintf("%s-%s.%s", repoSlug, name, r.FileExtension())
		}

		dest := filepath.Join(outDir, filename)
		if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dest, err)
		}
		log.Printf("wrote %s", dest)
	}
	return nil
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	defer func() { _ = f.Close() }() // os.File write errors surface via Encode

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	return nil
}
