// Copyright 2024 BlueFunda, Inc.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchReposByTopic(t *testing.T) {
	tests := []struct {
		name       string
		org        string
		topic      string
		pages      [][]string // repo names per page
		wantNames  []string
		wantErrSub string
	}{
		{
			name:      "single page",
			org:       "acme",
			topic:     "active",
			pages:     [][]string{{"alpha", "beta", "gamma"}},
			wantNames: []string{"alpha", "beta", "gamma"},
		},
		{
			name:      "empty result",
			org:       "acme",
			topic:     "no-such-topic",
			pages:     [][]string{{}},
			wantNames: nil,
		},
		{
			name:         "api error",
			org:          "acme",
			topic:        "active",
			wantErrSub:   "search repos by topic",
			// pages is nil — server will return 500
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pageIdx := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.pages == nil {
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}
				names := tc.pages[pageIdx]
				pageIdx++

				type item struct {
					Name string `json:"name"`
				}
				type result struct {
					Items []item `json:"items"`
				}
				items := make([]item, len(names))
				for i, n := range names {
					items[i] = item{Name: n}
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(result{Items: items})
			}))
			defer srv.Close()

			c := NewClient("test-token")
			c.baseURL = srv.URL

			got, err := c.SearchReposByTopic(tc.org, tc.topic)

			if tc.wantErrSub != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrSub)
				}
				if msg := err.Error(); len(msg) == 0 {
					t.Fatalf("expected error containing %q, got empty string", tc.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.wantNames) {
				t.Fatalf("got %d names, want %d: %v", len(got), len(tc.wantNames), got)
			}
			for i, name := range got {
				if name != tc.wantNames[i] {
					t.Errorf("name[%d] = %q, want %q", i, name, tc.wantNames[i])
				}
			}
		})
	}
}
