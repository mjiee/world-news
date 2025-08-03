package textx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mjiee/gokit/slicex"
)

// CleanText clean text
func CleanText(text string) string {
	if text == "" {
		return ""
	}

	// remove extra space
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// remove html tag
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")

	// replace special char
	text = strings.ReplaceAll(text, "\u00a0", " ") // replace non-breaking space
	text = strings.ReplaceAll(text, "\u200b", "")  // replace zero width space
	text = strings.ReplaceAll(text, "\ufeff", "")  // replace byte order mark

	return text
}

// MatchesKeyword checks if the link's title matches any of the keywords.
func MatchesKeyword(text string, keywords []string) (string, bool) {
	text = strings.ToLower(strings.TrimSpace(text))

	keyword := slicex.Find(keywords, func(k string) bool {
		keywordLower := strings.ToLower(k)

		if text == keywordLower {
			return true
		}

		if strings.Contains(text, keywordLower) {
			return true
		}

		pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(keywordLower))
		matched, _ := regexp.MatchString(pattern, text)
		if matched {
			return true
		}

		return false
	})

	return keyword, keyword != ""
}

// SimilarText check if two text are similar
func SimilarText(a, b string) bool {
	var (
		minLen   = min(len(a), len(b))
		maxLen   = max(len(a), len(b))
		matchLen = 0
	)

	for idx := range a[:minLen] {
		if a[idx] == b[idx] {
			matchLen++
		} else {
			break
		}
	}

	return (float64(matchLen) / float64(maxLen)) >= 0.9
}
