package valueobject

import "encoding/json"

// NewsWebsite represents a news website.
type NewsWebsite struct {
	Url   string   `json:"url"`             // url
	Attrs []string `json:"attrs,omitempty"` // html attributes
}

var (
	// NewsWebsiteCollection represents a collection of news websites.
	NewsWebsiteCollection = []NewsWebsite{
		{
			Url:   "https://www.world-newspapers.com/",
			Attrs: []string{"a.country-link, a.magazine-link", "a.magazine-link"},
		},
	}
)

// NewNewsWebsites creates a new collection of news websites from the provided data.
func NewNewsWebsites(data string) []*NewsWebsite {
	if len(data) == 0 {
		return nil
	}

	var websites []*NewsWebsite

	_ = json.Unmarshal([]byte(data), &websites)

	return websites
}
