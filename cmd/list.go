package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/railgun-0402/ops-changelog/internal/formatter"
	gh "github.com/railgun-0402/ops-changelog/internal/github"
)

var (
	flagRepo    string
	flagService string
	flagSince   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List merged PRs for a service",
	Example: `  ops-changelog list --repo myorg/myrepo --service apply-api --since 24h
  ops-changelog list --repo myorg/myrepo --service apply-api --since 7d`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVar(&flagRepo, "repo", "", "GitHub repository in owner/repo format (required)")
	listCmd.Flags().StringVar(&flagService, "service", "", "Service name — matches label 'service:<name>' (optional)")
	listCmd.Flags().StringVar(&flagSince, "since", "24h", "How far back to look (e.g. 1h, 24h, 7d, 30d)")

	_ = listCmd.MarkFlagRequired("repo")

	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, _ []string) error {
	owner, repo, err := splitRepo(flagRepo)
	if err != nil {
		return err
	}

	since, err := parseSince(flagSince)
	if err != nil {
		return fmt.Errorf("invalid --since value %q: use formats like 1h, 24h, 7d, 30d", flagSince)
	}

	client, err := gh.NewClient(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return err
	}

	prs, err := client.ListMergedPRs(context.Background(), gh.ListMergedPRsOptions{
		Owner:   owner,
		Repo:    repo,
		Service: flagService,
		Since:   since,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch PRs: %w", err)
	}

	formatter.PrintPRs(os.Stdout, prs, flagRepo, flagService, since)
	return nil
}

func splitRepo(repo string) (owner, name string, err error) {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("--repo must be in owner/repo format, got %q", repo)
	}
	return parts[0], parts[1], nil
}

// parseSince supports Go duration strings (1h, 24h) and day shorthand (7d, 30d).
func parseSince(s string) (time.Time, error) {
	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return time.Time{}, err
		}
		return time.Now().UTC().AddDate(0, 0, -n), nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().UTC().Add(-d), nil
}
