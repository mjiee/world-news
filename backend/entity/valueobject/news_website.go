package valueobject

import (
	"github.com/mjiee/world-news/backend/pkg/urlx"
)

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
		{
			Url: "https://onlinenewspapers.com/",
			Selector: &Selector{
				Website: ".t3b", Child: &Selector{Website: ".t3b", Child: &Selector{Website: ".t3b",
					Child: &Selector{Website: ".t3b"}}},
			},
		},
	}
)

// GetHost returns the host of the news website.
func (nw *NewsWebsite) GetHost() string {
	return urlx.ExtractHostFromURL(nw.Url)
}
