package formatter

import (
	"encoding/json"
	"io"

	gh "github.com/railgun-0402/ops-changelog/internal/github"
)

func PrintPRsJSON(w io.Writer, prs []gh.PR) error {
	return json.NewEncoder(w).Encode(prs)
}
