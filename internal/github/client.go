package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const baseURL = "https://api.github.com"

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

// apiPR is the raw GitHub API pull request shape (subset of fields we need).
type apiPR struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	HTMLUrl  string `json:"html_url"`
	MergedAt string `json:"merged_at"`
	User     struct {
		Login string `json:"login"`
	} `json:"user"`
	Labels []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

// fetchMergedPRPage fetches one page of merged PRs (sorted by updated desc).
// Returns the PRs, whether there are more pages, and any error.
func (c *Client) fetchMergedPRPage(ctx context.Context, owner, repo string, page int) ([]PR, bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=closed&sort=updated&direction=desc&per_page=100&page=%d",
		baseURL, owner, repo, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return nil, false, fmt.Errorf("GitHub API returned %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var raw []apiPR
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, false, err
	}

	var prs []PR
	for _, r := range raw {
		// Skip unmerged closed PRs
		if r.MergedAt == "" {
			continue
		}
		mergedAt, err := time.Parse(time.RFC3339, r.MergedAt)
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

	hasMore := len(raw) == 100
	return prs, hasMore, nil
}
