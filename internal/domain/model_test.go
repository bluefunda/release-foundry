package domain

import "testing"

func TestClassifyPR(t *testing.T) {
	tests := []struct {
		name   string
		labels []string
		title  string
		want   string
	}{
		{"feature label wins", []string{"fix", "feature"}, "", "feature"},
		{"fix label only", []string{"fix"}, "", "fix"},
		{"performance label", []string{"performance"}, "", "performance"},
		{"security label", []string{"security"}, "", "security"},
		{"infrastructure label", []string{"infrastructure"}, "", "infrastructure"},
		{"no matching label", []string{"docs"}, "", "other"},
		{"empty labels no prefix", []string{}, "random title", "other"},
		{"feature beats performance", []string{"performance", "feature"}, "", "feature"},
		// Title prefix fallback.
		{"feat prefix", []string{}, "feat: add login", "feature"},
		{"fix prefix", []string{}, "fix: broken button", "fix"},
		{"perf prefix", []string{}, "perf: speed up queries", "performance"},
		{"feat with scope", []string{}, "feat(auth): add SSO", "feature"},
		// Labels take priority over title.
		{"label overrides title", []string{"fix"}, "feat: something", "fix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyPR(tt.labels, tt.title)
			if got != tt.want {
				t.Errorf("ClassifyPR(%v, %q) = %q, want %q", tt.labels, tt.title, got, tt.want)
			}
		})
	}
}

func TestInferTypeFromTitle(t *testing.T) {
	tests := []struct {
		title string
		want  string
		ok    bool
	}{
		{"feat: add login", "feature", true},
		{"fix: broken button", "fix", true},
		{"perf: speed up queries", "performance", true},
		{"feat(auth): add SSO", "feature", true},
		{"chore: update deps", "chore", false},
		{"random title", "random title", false},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got, ok := InferTypeFromTitle(tt.title)
			if got != tt.want || ok != tt.ok {
				t.Errorf("InferTypeFromTitle(%q) = (%q, %v), want (%q, %v)", tt.title, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestIsExcludedByTitle(t *testing.T) {
	tests := []struct {
		title string
		want  bool
	}{
		{"chore: update deps", true},
		{"refactor: clean up", true},
		{"docs: add readme", true},
		{"feat: new feature", false},
		{"fix: bug fix", false},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if got := IsExcludedByTitle(tt.title); got != tt.want {
				t.Errorf("IsExcludedByTitle(%q) = %v, want %v", tt.title, got, tt.want)
			}
		})
	}
}
