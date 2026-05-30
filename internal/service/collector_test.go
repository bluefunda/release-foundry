package service

import "testing"

func TestExtractSection_CustomerImpact(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "h2 section",
			body: "## Summary\nSome stuff\n## Customer Impact\nUsers can now export CSVs.\n## Notes\nDone.",
			want: "Users can now export CSVs.",
		},
		{
			name: "h3 section",
			body: "### Customer Impact\nFaster page loads.\n### Other\nNope.",
			want: "Faster page loads.",
		},
		{
			name: "no section",
			body: "Just a normal PR body without any special section.",
			want: "",
		},
		{
			name: "section at end of body",
			body: "## Summary\nStuff\n## Customer Impact\nReduced latency by 50%.",
			want: "Reduced latency by 50%.",
		},
		{
			name: "multiline content",
			body: "## Customer Impact\nLine one.\nLine two.\n## Next",
			want: "Line one.\nLine two.",
		},
		{
			name: "crlf line endings",
			body: normalizeBody("## Customer Impact\r\nWindows line endings.\r\n## Next"),
			want: "Windows line endings.",
		},
		{
			name: "empty section",
			body: "## Customer Impact\n\n## Next",
			want: "",
		},
		{
			name: "header with trailing text",
			body: "## Customer Impact (Required for Feature / Perf / Security)\nUsers benefit from SSO.\n## Test Plan",
			want: "Users benefit from SSO.",
		},
		{
			name: "html comments stripped",
			body: "## Customer Impact\n<!-- Who benefits? -->\nFaster deployments.\n## Next",
			want: "Faster deployments.",
		},
		{
			name: "only html comments",
			body: "## Customer Impact\n<!-- placeholder -->\n## Next",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSection(tt.body, "customer impact")
			if got != tt.want {
				t.Errorf("extractSection(customer impact) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractSection_Metrics(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "metrics present",
			body: "## Summary\nStuff\n## Metrics\n- Latency change: -40ms p99\n## Next",
			want: "- Latency change: -40ms p99",
		},
		{
			name: "metrics absent",
			body: "## Summary\nStuff\n## Customer Impact\nNone.",
			want: "",
		},
		{
			name: "metrics with comments stripped",
			body: "## Metrics\n<!-- Optional -->\n- Cost reduction: 30%\n## Next",
			want: "- Cost reduction: 30%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSection(tt.body, "metrics")
			if got != tt.want {
				t.Errorf("extractSection(metrics) = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractSection_MarketingNotes(t *testing.T) {
	body := "## Marketing Notes\nBig launch item for Q1.\n## Test Plan\n- [x] done"
	got := extractSection(body, "marketing notes")
	want := "Big launch item for Q1."
	if got != want {
		t.Errorf("extractSection(marketing notes) = %q, want %q", got, want)
	}
}

func TestExtractSummary(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "has summary section",
			body: "## Summary\nAdded SSO login.\n## Customer Impact\nUsers benefit.",
			want: "Added SSO login.",
		},
		{
			name: "legacy what is this change",
			body: "## What is this change?\nRefactored auth module.\n## Type\n- [x] feature",
			want: "Refactored auth module.",
		},
		{
			name: "no recognized section falls back to full body",
			body: "Just some text about the PR.",
			want: "Just some text about the PR.",
		},
		{
			name: "summary with html comments",
			body: "## Summary\n<!-- Brief description -->\nNew caching layer.\n## Type",
			want: "New caching layer.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSummary(tt.body)
			if got != tt.want {
				t.Errorf("extractSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}
