package valueobject

import "net/url"

// NewsWebsite represents a news website.
type NewsWebsite struct {
	Url      string    `json:"url"`                // url
	Selector *Selector `json:"selector,omitempty"` // selectors
}

// Selector represents a selector.
type Selector struct {
	// News website selector
	Website string `json:"website,omitempty"`

	// News selector
	Topic string `json:"topic,omitempty"`
	Link  string `json:"link,omitempty"`

	// News detail selector
	Title   string `json:"title,omitempty"`
	Time    string `json:"time,omitempty"`
	Image   string `json:"image,omitempty"`
	Content string `json:"content,omitempty"`
	Author  string `json:"author,omitempty"`

	Child *Selector `json:"child,omitempty"` // child selector
}

var (
	// NewsWebsiteCollection represents a collection of news websites.
	NewsWebsiteCollection = []NewsWebsite{
		{
			Url:      "https://www.world-newspapers.com/",
			Selector: &Selector{Website: "a.country-link, a.magazine-link", Child: &Selector{Website: "a.magazine-link"}},
		},
	}
)

// GetHost returns the host of the news website.
func (nw *NewsWebsite) GetHost() string {
	data, err := url.Parse(nw.Url)
	if err != nil {
		return nw.Url
	}

	return data.Host
}
