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
			name: "empty fields PR",
			prs: []gh.PR{
				{
					Number:   0,
					Title:    "",
					URL:      "",
					Author:   "",
					MergedAt: time.Time{},
					Labels:   []string{},
				},
			},
		},
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
			name: "multiple PRs",
			prs: []gh.PR{
				{
					Number:   1,
					Title:    "Fix something",
					URL:      "https://github.com/myorg/myrepo/pull/1",
					Author:   "alice",
					MergedAt: time.Date(2026, 3, 14, 10, 30, 0, 0, time.UTC),
					Labels:   []string{"service:apply-api", "fix"},
				},
				{
					Number:   2,
					Title:    "Fix something2",
					URL:      "https://github.com/myorg/myrepo/pull/2",
					Author:   "Bob",
					MergedAt: time.Date(2026, 3, 30, 22, 15, 0, 0, time.UTC),
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

			for i := range tt.prs {
				if got[i].Number != tt.prs[i].Number {
					t.Errorf("PR[%d] Number: got %d, want %d", i, got[i].Number, tt.prs[i].Number)
				}
				if got[i].Title != tt.prs[i].Title {
					t.Errorf("PR[%d] Title: got %q, want %q", i, got[i].Title, tt.prs[i].Title)
				}
				if got[i].Author != tt.prs[i].Author {
					t.Errorf("PR[%d] Author: got %q, want %q", i, got[i].Author, tt.prs[i].Author)
				}
			}
		})
	}
}
