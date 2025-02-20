package valueobject

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
