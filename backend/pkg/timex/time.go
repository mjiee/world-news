package timex

import (
	"regexp"
	"time"

	"github.com/mjiee/world-news/backend/pkg/textx"

	"github.com/markusmobius/go-dateparser"
)

var (
	urlDateRegex   = regexp.MustCompile(`(?:\d{4}[-/]?\d{2}[-/]?\d{2}|\d{2}[-/]?\d{2}[-/]?\d{4})(?:[T\s]\d{2}:\d{2}:\d{2})?`)
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

	match := urlDateRegex.FindString(text)
	date, err = dateparser.Parse(nil, match)
	if err == nil {
		return date.Time
	}

	matchs := timestampRegex.FindStringSubmatch(text)
	if len(matchs) == 0 {
		return time.Time{}
	}

	date, _ = dateparser.Parse(nil, matchs[1])

	return date.Time
}
