package formatter

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	gh "github.com/railgun-0402/ops-changelog/internal/github"
)

func TestPrintPRsJSON(t *testing.T) {
	tests := []struct {
		name string
		prs  []gh.PR
	}{
		{
			name: "single PR",
			prs: []gh.PR{
				{
					Number:   1,
					Title:    "Fix something",
					URL:      "https://github.com/myorg/myrepo/pull/1",
					Author:   "alice",
					MergedAt: time.Date(2026, 3, 14, 10, 30, 0, 0, time.UTC),
					Labels:   []string{"service:apply-api", "fix"},
				},
			},
		},
		{
			name: "zero PRs",
			prs:  []gh.PR{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := PrintPRsJSON(&buf, tt.prs); err != nil {
				t.Fatalf("PrintPRsJSON returned error: %v", err)
			}

			var got []gh.PR
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			if len(got) != len(tt.prs) {
				t.Fatalf("got %d PRs, want %d", len(got), len(tt.prs))
			}

			// Check the zero-value PR is preserved
			if len(got) == 0 {
				return
			}

			if got[0].Number != tt.prs[0].Number {
				t.Errorf("Number: got %d, want %d", got[0].Number, tt.prs[0].Number)
			}
			if got[0].Title != tt.prs[0].Title {
				t.Errorf("Title: got %q, want %q", got[0].Title, tt.prs[0].Title)
			}
			if got[0].Author != tt.prs[0].Author {
				t.Errorf("Author: got %q, want %q", got[0].Author, tt.prs[0].Author)
			}
		})
	}
}
