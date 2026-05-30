package domain

import (
	"strings"
	"time"
)

// PullRequest holds the extracted fields from a merged GitHub PR.
type PullRequest struct {
	Number           int      `json:"number"`
	Type             string   `json:"type"`
	Title            string   `json:"title"`
	CustomerImpact   string   `json:"customer_impact_raw"`
	TechnicalSummary string   `json:"technical_summary"`
	Metrics          string   `json:"metrics,omitempty"`
	MarketingNotes   string   `json:"marketing_notes,omitempty"`
	Labels           []string `json:"labels"`
	Author           string   `json:"author"`
	FilesChanged     int      `json:"files_changed"`
	Additions        int      `json:"additions"`
	Deletions        int      `json:"deletions"`
	MergedAt         string   `json:"merged_at"`
}

// SummaryStats holds aggregate counts by PR type.
type SummaryStats struct {
	TotalPRs    int `json:"total_prs"`
	Features    int `json:"features"`
	Fixes       int `json:"fixes"`
	Performance int `json:"performance"`
}

// WeeklySummary is the top-level output structure written to JSON.
type WeeklySummary struct {
	GeneratedAt    string        `json:"generated_at"`
	Repository     string        `json:"repository"`
	Edition        string        `json:"edition,omitempty"`
	TimeWindowDays int           `json:"time_window_days"`
	SummaryStats   SummaryStats  `json:"summary_stats"`
	PullRequests   []PullRequest `json:"pull_requests"`
}

// BatchSummary is the combined output when processing multiple repos.
type BatchSummary struct {
	GeneratedAt    string          `json:"generated_at"`
	TimeWindowDays int             `json:"time_window_days"`
	Repositories   []WeeklySummary `json:"repositories"`
}

// IncludeLabels are the PR labels we want to collect.
var IncludeLabels = map[string]bool{
	"feature":        true,
	"fix":            true,
	"performance":    true,
	"security":       true,
	"infrastructure": true,
}

// ExcludeLabels are the PR labels we must skip.
var ExcludeLabels = map[string]bool{
	"internal": true,
	"refactor": true,
	"chore":    true,
}

// TitlePrefixMap maps conventional commit prefixes to canonical types.
var TitlePrefixMap = map[string]string{
	"feat":     "feature",
	"fix":      "fix",
	"perf":     "performance",
	"security": "security",
	"infra":    "infrastructure",
}

// TitleExcludePrefixes are conventional commit prefixes that map to excluded types.
var TitleExcludePrefixes = map[string]bool{
	"chore":    true,
	"refactor": true,
	"docs":     true,
	"ci":       true,
	"test":     true,
	"style":    true,
}

// InferTypeFromTitle extracts a conventional commit prefix (e.g. "feat:" → "feature").
// Returns the canonical type and whether it matched.
func InferTypeFromTitle(title string) (string, bool) {
	// Conventional format: "type: description" or "type(scope): description"
	prefix := strings.ToLower(title)
	for _, sep := range []string{"(", ":"} {
		if idx := strings.Index(prefix, sep); idx > 0 {
			prefix = prefix[:idx]
			break
		}
	}
	prefix = strings.TrimSpace(prefix)

	if mapped, ok := TitlePrefixMap[prefix]; ok {
		return mapped, true
	}
	return prefix, false
}

// IsExcludedByTitle returns true if the title's conventional prefix maps to an excluded type.
func IsExcludedByTitle(title string) bool {
	prefix := strings.ToLower(title)
	for _, sep := range []string{"(", ":"} {
		if idx := strings.Index(prefix, sep); idx > 0 {
			prefix = prefix[:idx]
			break
		}
	}
	return TitleExcludePrefixes[strings.TrimSpace(prefix)]
}

// ClassifyPR determines the primary type from labels, falling back to title prefix.
// Priority: feature > fix > performance > security > infrastructure.
func ClassifyPR(labels []string, title string) string {
	priority := []string{"feature", "fix", "performance", "security", "infrastructure"}
	set := make(map[string]bool, len(labels))
	for _, l := range labels {
		set[l] = true
	}
	for _, p := range priority {
		if set[p] {
			return p
		}
	}
	// Fallback: infer from conventional commit title prefix.
	if t, ok := InferTypeFromTitle(title); ok {
		return t
	}
	return "other"
}

// Config holds runtime configuration for the collector.
type Config struct {
	Token      string
	Owner      string
	Repo       string
	BaseBranch string
	WindowDays int
	Since      time.Time
	Edition    string
	Filters    FilterConfig
}
