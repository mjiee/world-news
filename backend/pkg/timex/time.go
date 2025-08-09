package timex

import (
	"regexp"
	"strings"
	"time"

	"github.com/mjiee/world-news/backend/pkg/textx"

	"github.com/markusmobius/go-dateparser"
)

var (
	urlDateRegex   = regexp.MustCompile(`(?:^|/)(\d{4})(?:[-/](\d{1,2})(?:[-/](\d{1,2}))?)?(?:/|$)`)
	timestampRegex = regexp.MustCompile(`(?:^|[^\d])(\d{10}|\d{13})(?:[^\d]|$)`)
)

// ParseTime
func ParseTime(text string) time.Time {
	text = textx.CleanText(text)

	if text == "" {
		return time.Time{}
	}

	date, err := dateparser.Parse(nil, text)
	if err == nil {
		return date.Time
	}

	dateStr := extractDateFromUrl(text)
	if dateStr == "" {
		return time.Time{}
	}

	date, _ = dateparser.Parse(nil, dateStr)

	return date.Time
}

// extractDateFromUrl extract date from url
func extractDateFromUrl(url string) string {
	matches := urlDateRegex.FindStringSubmatch(url)

	if len(matches) > 1 {
		parts := []string{matches[1]}
		if matches[2] != "" {
			parts = append(parts, matches[2])
			if matches[3] != "" {
				parts = append(parts, matches[3])
			}
		}

		return strings.Join(parts, "/")
	}

	matchs := timestampRegex.FindStringSubmatch(url)
	if len(matchs) >= 1 {
		return matchs[1]
	}

	return ""
}
