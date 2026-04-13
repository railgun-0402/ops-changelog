package github

import (
	"context"
	"fmt"
	"time"
)

// PR represents a merged pull request.
type PR struct {
	Number   int       `json:"number"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Author   string    `json:"author"`
	MergedAt time.Time `json:"merged_at"`
	Labels   []string  `json:"labels"`
}

// ListMergedPRsOptions holds filtering options.
type ListMergedPRsOptions struct {
	Owner   string
	Repo    string
	Service string // matches label "service:<Service>"
	Since   time.Time
}

// ListMergedPRs fetches merged PRs filtered by service label and merged-after time.
// It uses the GitHub Search API to filter server-side, avoiding full pagination
// through all closed PRs.
func (c *Client) ListMergedPRs(ctx context.Context, opts ListMergedPRsOptions) ([]PR, error) {
	var label string
	if opts.Service != "" {
		label = fmt.Sprintf("service:%s", opts.Service)
	}

	var result []PR
	page := 1

	for {
		prs, hasMore, err := c.SearchMergedPRsPage(ctx, opts.Owner, opts.Repo, label, opts.Since, page)
		if err != nil {
			return nil, err
		}

		result = append(result, prs...)

		if !hasMore {
			break
		}
		page++
	}

	return result, nil
}
