package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const searchURL = "https://api.github.com/search/issues"

// Client is a minimal GitHub REST API client.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a Client, resolving the token from GITHUB_TOKEN env or `gh auth token`.
func NewClient(token string) (*Client, error) {
	if token == "" {
		var err error
		token, err = tokenFromGHCLI()
		if err != nil {
			return nil, fmt.Errorf("no GitHub token: set GITHUB_TOKEN or run `gh auth login`")
		}
	}
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func tokenFromGHCLI() (string, error) {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// searchIssuesResponse is the GitHub Search Issues API response shape.
type searchIssuesResponse struct {
	TotalCount int             `json:"total_count"`
	Items      []searchIssuePR `json:"items"`
}

// searchIssuePR is one item from the Search Issues API (PR subset).
type searchIssuePR struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HTMLUrl string `json:"html_url"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
	PullRequest struct {
		MergedAt string `json:"merged_at"`
	} `json:"pull_request"`
}

// SearchMergedPRsPage fetches one page of merged PRs via GitHub Search API.
// Filtering by repo, optional label, and merged-after time is done server-side,
// which avoids paginating through unrelated closed PRs.
func (c *Client) SearchMergedPRsPage(ctx context.Context, owner, repo, label string, since time.Time, page int) ([]PR, bool, error) {
	q := fmt.Sprintf("is:pr is:merged repo:%s/%s merged:>=%s",
		owner, repo, since.UTC().Format(time.RFC3339))
	if label != "" {
		q += " label:" + label
	}

	rawURL := buildSearchURL(q, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub Search API returned %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var raw searchIssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, false, err
	}

	prs := make([]PR, 0, len(raw.Items))
	for _, r := range raw.Items {
		mergedAt, err := time.Parse(time.RFC3339, r.PullRequest.MergedAt)
		if err != nil {
			continue
		}
		labels := make([]string, len(r.Labels))
		for i, l := range r.Labels {
			labels[i] = l.Name
		}
		prs = append(prs, PR{
			Number:   r.Number,
			Title:    r.Title,
			URL:      r.HTMLUrl,
			Author:   r.User.Login,
			MergedAt: mergedAt,
			Labels:   labels,
		})
	}

	hasMore := len(raw.Items) == 100
	return prs, hasMore, nil
}

// buildSearchURL constructs the search API URL with properly encoded query parameters.
func buildSearchURL(q string, page int) string {
	params := url.Values{}
	params.Set("q", q)
	params.Set("sort", "created")
	params.Set("direction", "desc")
	params.Set("per_page", "100")
	params.Set("page", fmt.Sprintf("%d", page))
	return searchURL + "?" + params.Encode()
}
