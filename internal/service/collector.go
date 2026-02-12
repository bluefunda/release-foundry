package service

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/release-foundry/internal/domain"
	gh "github.com/release-foundry/internal/github"
)

// Collector orchestrates fetching, filtering, and structuring PR data.
type Collector struct {
	client *gh.Client
	cfg    domain.Config
}

// NewCollector creates a Collector with the given GitHub client and config.
func NewCollector(client *gh.Client, cfg domain.Config) *Collector {
	return &Collector{client: client, cfg: cfg}
}

// Collect fetches PRs, applies filtering rules, enriches with detail, and returns a WeeklySummary.
func (c *Collector) Collect() (*domain.WeeklySummary, error) {
	log.Printf("fetching merged PRs for %s/%s (base=%s, since=%s)",
		c.cfg.Owner, c.cfg.Repo, c.cfg.BaseBranch, c.cfg.Since.Format(time.RFC3339))

	items, err := c.client.ListMergedPRs(c.cfg.Owner, c.cfg.Repo)
	if err != nil {
		return nil, fmt.Errorf("list merged PRs: %w", err)
	}
	log.Printf("fetched %d closed PRs from API", len(items))

	var prs []domain.PullRequest
	stats := domain.SummaryStats{}

	for _, item := range items {
		// Must be actually merged (not just closed).
		if item.MergedAt == "" {
			continue
		}

		mergedAt, err := time.Parse(time.RFC3339, item.MergedAt)
		if err != nil {
			log.Printf("skip PR #%d: cannot parse merged_at %q: %v", item.Number, item.MergedAt, err)
			continue
		}

		// Time window filter.
		if mergedAt.Before(c.cfg.Since) {
			continue
		}

		// Base branch filter.
		if item.Base.Ref != c.cfg.BaseBranch {
			continue
		}

		labels := extractLabelNames(item.Labels)

		// Exclude: by label first, then by title prefix.
		if hasAnyLabel(labels, domain.ExcludeLabels) {
			log.Printf("skip PR #%d %q: has excluded label", item.Number, item.Title)
			continue
		}
		if len(labels) == 0 && domain.IsExcludedByTitle(item.Title) {
			log.Printf("skip PR #%d %q: excluded by title prefix", item.Number, item.Title)
			continue
		}

		// Include: by label, or by title prefix when no labels exist.
		hasIncludeLabel := hasAnyLabel(labels, domain.IncludeLabels)
		_, titleMatched := domain.InferTypeFromTitle(item.Title)
		if !hasIncludeLabel && !(len(labels) == 0 && titleMatched) {
			log.Printf("skip PR #%d %q: no include label or recognized title prefix", item.Number, item.Title)
			continue
		}

		// Fetch file-change stats.
		detail, err := c.client.GetPRDetail(c.cfg.Owner, c.cfg.Repo, item.Number)
		if err != nil {
			log.Printf("warning: could not fetch detail for PR #%d: %v (skipping stats)", item.Number, err)
			detail = &gh.PRDetail{}
		}

		prType := domain.ClassifyPR(labels, item.Title)

		body := normalizeBody(item.Body)

		pr := domain.PullRequest{
			Number:           item.Number,
			Type:             prType,
			Title:            item.Title,
			CustomerImpact:   extractSection(body, "customer impact"),
			TechnicalSummary: extractSummary(body),
			Metrics:          extractSection(body, "metrics"),
			MarketingNotes:   extractSection(body, "marketing notes"),
			Labels:           labels,
			Author:           item.User.Login,
			FilesChanged:     detail.ChangedFiles,
			Additions:        detail.Additions,
			Deletions:        detail.Deletions,
			MergedAt:         item.MergedAt,
		}

		prs = append(prs, pr)

		// Tally stats.
		stats.TotalPRs++
		switch prType {
		case "feature":
			stats.Features++
		case "fix":
			stats.Fixes++
		case "performance":
			stats.Performance++
		}
	}

	log.Printf("collected %d PRs after filtering", len(prs))

	summary := &domain.WeeklySummary{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		Repository:     fmt.Sprintf("%s/%s", c.cfg.Owner, c.cfg.Repo),
		TimeWindowDays: c.cfg.WindowDays,
		SummaryStats:   stats,
		PullRequests:   prs,
	}

	// Ensure pull_requests is always an array, never null.
	if summary.PullRequests == nil {
		summary.PullRequests = []domain.PullRequest{}
	}

	return summary, nil
}

func extractLabelNames(labels []gh.Label) []string {
	names := make([]string, 0, len(labels))
	for _, l := range labels {
		names = append(names, l.Name)
	}
	return names
}

func hasAnyLabel(labels []string, set map[string]bool) bool {
	for _, l := range labels {
		if set[l] {
			return true
		}
	}
	return false
}

func normalizeBody(body string) string {
	return strings.ReplaceAll(body, "\r\n", "\n")
}

// sectionRe builds a regex that matches a markdown heading containing the given name,
// capturing everything until the next heading or end of string.
// Tolerates trailing text after the section name (e.g. "## Customer Impact (Required)").
func sectionRe(name string) *regexp.Regexp {
	// Escape the name for regex, then allow flexible whitespace between words.
	words := strings.Fields(name)
	for i, w := range words {
		words[i] = regexp.QuoteMeta(w)
	}
	pattern := fmt.Sprintf(`(?i)(?:^|\n)#{1,3}[^\S\n]*%s[^\n]*\n([\s\S]*?)(?:\n#{1,3}\s|\z)`, strings.Join(words, `\s+`))
	return regexp.MustCompile(pattern)
}

// extractSection pulls the body text under a markdown heading that contains `name`.
// HTML comments (<!-- ... -->) are stripped from the extracted content.
func extractSection(body, name string) string {
	re := sectionRe(name)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	return stripComments(strings.TrimSpace(matches[1]))
}

// extractSummary returns the ## Summary section if present, otherwise the full body.
func extractSummary(body string) string {
	if s := extractSection(body, "summary"); s != "" {
		return s
	}
	// Legacy fallback: "What is this change?"
	if s := extractSection(body, "what is this change"); s != "" {
		return s
	}
	return stripComments(strings.TrimSpace(body))
}

var htmlCommentRe = regexp.MustCompile(`<!--[\s\S]*?-->`)

func stripComments(s string) string {
	return strings.TrimSpace(htmlCommentRe.ReplaceAllString(s, ""))
}
