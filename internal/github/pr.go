package github

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// PR represents a merged pull request.
type PR struct {
	Number   int
	Title    string
	URL      string
	Author   string
	MergedAt time.Time
	Labels   []string
}

// ListMergedPRsOptions holds filtering options.
type ListMergedPRsOptions struct {
	Owner   string
	Repo    string
	Service string // matches label "service:<Service>"
	Since   time.Time
}

// ListMergedPRs fetches merged PRs filtered by service label and merged-after time.
func (c *Client) ListMergedPRs(ctx context.Context, opts ListMergedPRsOptions) ([]PR, error) {
	targetLabel := fmt.Sprintf("service:%s", opts.Service)

	var result []PR
	page := 1

	for {
		prs, hasMore, err := c.fetchMergedPRPage(ctx, opts.Owner, opts.Repo, page)
		if err != nil {
			return nil, err
		}

		for _, pr := range prs {
			// Stop pagination once we go past the since boundary
			if pr.MergedAt.Before(opts.Since) {
				return result, nil
			}

			if hasLabel(pr.Labels, targetLabel) {
				result = append(result, pr)
			}
		}

		if !hasMore {
			break
		}
		page++
	}

	return result, nil
}

func hasLabel(labels []string, target string) bool {
	for _, l := range labels {
		if strings.EqualFold(l, target) {
			return true
		}
	}
	return false
}
