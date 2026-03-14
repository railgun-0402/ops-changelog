package formatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	gh "ops-changelog/internal/github"
)

// PrintPRs writes a human-readable PR list to w.
func PrintPRs(w io.Writer, prs []gh.PR, repo, service string, since time.Time) {
	serviceLabel := "all services"
	if service != "" {
		serviceLabel = "service:" + service
	}

	if len(prs) == 0 {
		fmt.Fprintf(w, "No merged PRs found for %s in %s since %s\n",
			serviceLabel, repo, since.Format("2006-01-02 15:04 UTC"))
		return
	}

	fmt.Fprintf(w, "Merged PRs for %s in %s (since %s)\n",
		serviceLabel, repo, since.Format("2006-01-02 15:04 UTC"))
	fmt.Fprintln(w, strings.Repeat("─", 72))

	for _, pr := range prs {
		fmt.Fprintf(w, "[%s]  #%d  %s\n",
			pr.MergedAt.UTC().Format("2006-01-02 15:04"), pr.Number, pr.Title)
		fmt.Fprintf(w, "  Author : %s\n", pr.Author)
		fmt.Fprintf(w, "  URL    : %s\n", pr.URL)
		if len(pr.Labels) > 0 {
			fmt.Fprintf(w, "  Labels : %s\n", strings.Join(pr.Labels, ", "))
		}
		fmt.Fprintln(w)
	}
}
