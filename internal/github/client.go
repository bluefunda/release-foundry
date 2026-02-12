package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.github.com"

// Client wraps authenticated access to the GitHub REST API.
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient creates a Client with the given personal access token.
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}
}

// Label represents a GitHub label.
type Label struct {
	Name string `json:"name"`
}

// PRListItem represents the subset of fields returned by the list-pulls endpoint.
type PRListItem struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	State    string `json:"state"`
	MergedAt string `json:"merged_at"`
	Base     struct {
		Ref string `json:"ref"`
	} `json:"base"`
	Labels []Label `json:"labels"`
	User   struct {
		Login string `json:"login"`
	} `json:"user"`
}

// PRDetail holds additional detail fetched per-PR.
type PRDetail struct {
	ChangedFiles int `json:"changed_files"`
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
}

// ListMergedPRs fetches all merged PRs for the repo, paginating through all results.
// It returns PRs in the order the API provides (most-recently updated first).
func (c *Client) ListMergedPRs(owner, repo string) ([]PRListItem, error) {
	var all []PRListItem
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=closed&base=main&sort=updated&direction=desc&per_page=%d&page=%d",
			baseURL, owner, repo, perPage, page)

		body, headers, err := c.get(url)
		if err != nil {
			return nil, fmt.Errorf("list PRs page %d: %w", page, err)
		}

		var items []PRListItem
		if err := json.Unmarshal(body, &items); err != nil {
			return nil, fmt.Errorf("decode PRs page %d: %w", page, err)
		}

		all = append(all, items...)

		if !hasNextPage(headers) || len(items) < perPage {
			break
		}
		page++
	}

	return all, nil
}

// GetPRDetail fetches the full PR object to obtain file-change stats.
func (c *Client) GetPRDetail(owner, repo string, number int) (*PRDetail, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", baseURL, owner, repo, number)
	body, _, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("get PR #%d detail: %w", number, err)
	}
	var detail PRDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("decode PR #%d detail: %w", number, err)
	}
	return &detail, nil
}

// get performs an authenticated GET and handles rate-limit backoff.
func (c *Client) get(url string) ([]byte, http.Header, error) {
	for {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		// Handle rate limiting.
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			if wait := rateLimitWait(resp.Header); wait > 0 {
				log.Printf("rate limited, waiting %s", wait)
				time.Sleep(wait)
				continue
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("read response body: %w", err)
		}

		if resp.StatusCode >= 400 {
			return nil, nil, fmt.Errorf("GitHub API %s returned %d: %s", url, resp.StatusCode, truncate(string(body), 200))
		}

		return body, resp.Header, nil
	}
}

// rateLimitWait returns how long to sleep based on X-RateLimit-Reset.
func rateLimitWait(h http.Header) time.Duration {
	reset := h.Get("X-RateLimit-Reset")
	if reset == "" {
		// Fallback from Retry-After header (secondary rate limits).
		if ra := h.Get("Retry-After"); ra != "" {
			sec, err := strconv.Atoi(ra)
			if err == nil {
				return time.Duration(sec) * time.Second
			}
		}
		return 5 * time.Second
	}
	epoch, err := strconv.ParseInt(reset, 10, 64)
	if err != nil {
		return 5 * time.Second
	}
	wait := time.Until(time.Unix(epoch, 0)) + 1*time.Second // +1s buffer
	if wait < 0 {
		return 0
	}
	return wait
}

// hasNextPage checks the Link header for a "next" relation.
func hasNextPage(h http.Header) bool {
	link := h.Get("Link")
	return strings.Contains(link, `rel="next"`)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
