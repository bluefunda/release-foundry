package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/release-foundry/internal/domain"
	gh "github.com/release-foundry/internal/github"
	"github.com/release-foundry/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token := flag.String("token", "", "GitHub personal access token (overrides GITHUB_TOKEN)")
	owner := flag.String("owner", "", "GitHub repository owner (overrides GITHUB_OWNER)")
	repo := flag.String("repo", "", "GitHub repository name (overrides GITHUB_REPO)")
	days := flag.Int("days", 7, "number of days to look back")
	output := flag.String("output", "weekly_engineering_summary.json", "output file path")
	flag.Parse()

	cfg, err := loadConfig(*token, *owner, *repo, *days)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	client := gh.NewClient(cfg.Token)
	collector := service.NewCollector(client, cfg)

	summary, err := collector.Collect()
	if err != nil {
		log.Fatalf("collection failed: %v", err)
	}

	if err := writeJSON(*output, summary); err != nil {
		log.Fatalf("write output: %v", err)
	}

	log.Printf("wrote %s (%d PRs)", *output, summary.SummaryStats.TotalPRs)
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

func loadConfig(flagToken, flagOwner, flagRepo string, days int) (domain.Config, error) {
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

	since := time.Now().UTC().AddDate(0, 0, -days)

	return domain.Config{
		Token:      token,
		Owner:      owner,
		Repo:       repo,
		BaseBranch: "main",
		WindowDays: days,
		Since:      since,
	}, nil
}

func writeJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	return nil
}
