package valueobject

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// NewsWebsite represents a news website.
type NewsWebsite struct {
	Url       string   `json:"url"`                 // url
	Selectors []string `json:"selectors,omitempty"` // selectors
}

var (
	// NewsWebsiteCollection represents a collection of news websites.
	NewsWebsiteCollection = []NewsWebsite{
		{
			Url:       "https://www.world-newspapers.com/",
			Selectors: []string{"a.country-link, a.magazine-link", "a.magazine-link"},
		},
	}
)

// NewsWebsitesFromAny converts any to NewsWebsite.
func NewsWebsitesFromAny(data any) ([]*NewsWebsite, error) {
	var websites []*NewsWebsite

	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(dataStr, &websites); err != nil {
		return nil, errors.WithStack(err)
	}

	return websites, nil
}
