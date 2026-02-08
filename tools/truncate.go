package tools

import "fmt"

func truncateWithNotice(content string, maxChars int) (string, bool) {
	if maxChars <= 0 || len(content) <= maxChars {
		return content, false
	}

	notice := fmt.Sprintf(
		"\n[Truncated] Output exceeded %d characters. Narrow the scope or read in chunks.",
		maxChars,
	)
	return content[:maxChars] + notice, true
}
