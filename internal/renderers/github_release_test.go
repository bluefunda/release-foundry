package renderers_test

import (
	"strings"
	"testing"

	"github.com/release-foundry/internal/domain"
	"github.com/release-foundry/internal/renderers"
)

func TestGithubRelease_groupsAndOrders(t *testing.T) {
	summary := domain.WeeklySummary{
		Repository: "bluefunda/btp-go",
		PullRequests: []domain.PullRequest{
			{Number: 10, Type: "fix", Title: "fix: nil pointer in connectivity", Author: "alice"},
			{Number: 11, Type: "feature", Title: "feat: add destination caching", Author: "bob"},
			{Number: 12, Type: "security", Title: "security: bump crypto dep", Author: "carol"},
		},
	}

	out := renderers.GithubRelease(summary)

	// Features must appear before fixes (typeOrder)
	featIdx := strings.Index(out, "✨ Features")
	fixIdx := strings.Index(out, "🐛 Bug Fixes")
	secIdx := strings.Index(out, "🔒 Security")

	if featIdx == -1 || fixIdx == -1 || secIdx == -1 {
		t.Fatalf("missing section headers:\n%s", out)
	}
	if featIdx > fixIdx {
		t.Errorf("Features section should precede Bug Fixes")
	}
	if fixIdx > secIdx {
		t.Errorf("Bug Fixes section should precede Security")
	}

	// Conventional prefix stripped from display title
	if !strings.Contains(out, "add destination caching") {
		t.Errorf("expected clean title without 'feat:' prefix")
	}
	if strings.Contains(out, "feat:") {
		t.Errorf("prefix 'feat:' should be stripped from display title")
	}

	// PR number and author included
	if !strings.Contains(out, "#11") || !strings.Contains(out, "@bob") {
		t.Errorf("PR number and author should appear in output")
	}
}

func TestGithubRelease_empty(t *testing.T) {
	out := renderers.GithubRelease(domain.WeeklySummary{Repository: "bluefunda/btp-go"})
	if !strings.Contains(out, "No user-facing changes") {
		t.Errorf("expected empty-state message, got:\n%s", out)
	}
}
