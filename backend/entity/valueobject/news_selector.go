package valueobject

import (
	"regexp"
	"slices"
	"strings"
	"time"
)

// maximum validity period
const MaxValidityPeriod = 7 * 24 * time.Hour

// common selectors
var (
	ExcludeSelectors = []string{
		"header", ".header", "#header", "#footer", "#footbot", ".footer", ".footer", "aside", ".aside", ".sidebar",
		".side-bar", "nav", ".nav", "navigation", ".navigation", ".menu", ".top-menu", ".main-menu",
		".ad", ".advertisement", ".ads", ".sponsor", ".promoted", ".banner", ".tag", ".tags", ".tag-list", ".categories",
		".popup", ".modal", ".overlay", ".newsletter", ".subscription", ".search", ".search-box", ".search-form",
		".comment", ".comments", ".comment-section", ".reply", ".replies", ".discuss",
		".social", ".social-share", ".share", ".sharing", ".follow", ".subscribe",
		".related", ".related-posts", ".related-articles", ".recommend", ".hot", ".popular", ".trending",
		".pagination", ".pager", ".next", ".prev", ".breadcrumb", ".breadcrumbs",
		".widget", ".plugin", ".embed", ".iframe", ".video-ad",
		"[class*='ad']", "[class*='advertisement']", "[class*='sponsor']",
		"[class*='comment']", "[class*='share']", "[class*='social']",
		"[class*='related']", "[class*='recommend']", "[class*='hot']",
		"[class*='sidebar']", "[class*='aside']", "[class*='footer']", "[class*='footbot']",
		"[class*='header']", "[class*='nav']", "[class*='menu']", "[class*='tag']",
		"script", "style", "noscript", "iframe", "embed", "object",
	}
	LinkSelector = "a[href]"
)

// news selectors
var (
	NewsTopicLinkSelectors = []string{
		"nav a, .nav a, .menu a, .navigation a, header a", "a[href]",
	}
	NewsItemSelectors = []string{
		"article", ".article", ".news-item", ".news-list-item",
		".story", ".story-item", ".post", ".post-item",
		".item", ".list-item", ".content-item", ".entry",
		".news", ".news-box", ".article-box", ".story-box",
		"li", ".li", "tr", ".row", ".card", ".box",
		"[class*='item']", "[class*='news']", "[class*='article']",
		"[class*='story']", "[class*='post']", "[class*='entry']",
		"[id*='item']", "[id*='news']", "[id*='article']",
	}
	NewsMainBodySelectors = []string{
		// structured selector
		"main", "article", "section[role='main']", "div[role='main']",

		// id selector
		"#main", "#main-content", "#primary", "#article", "#content", "#story", "#post", "#news-content",
		"#article-content", "#story-content", "#post-content", "#content-body", "#main-article", "#primary-content",
		"#page-content", "#site-content", "#entry-content",

		// class selector
		".wrapper article", ".page main", ".site-content article", ".content-area article", ".primary-content article",
		"h1 + div", "h1 + section", "h1 + article", ".headline + div", ".title + div",
		"div:has(p):not(:has(nav)):not(:has(.advertisement)):not(:has(.sidebar))",
		"section:has(p):not(:has(nav)):not(:has(.ads))",
		".story-container", ".article-container", ".news-container", ".content-container", ".text-container",
		".body-container", ".story-wrapper", ".article-wrapper", ".content-wrapper",
		"body > *:not(header):not(nav):not(footer):not(.sidebar):not(.advertisement)",
		".container > *:not(.header):not(.nav):not(.footer):not(.sidebar):not(.ads)",
		"main:not(.navigation):not(.sidebar):not(.footer)",
		"article:not(.related):not(.recommended):not(.advertisement)",
		".content:not(.sidebar-content):not(.footer-content):not(.nav-content)",
		".primary:not(.primary-nav):not(.primary-sidebar)", ".container > .row > .col:first-child .content",
		".wrapper > .main-content:first-of-type", ".page-container > .content-area > article",

		// class selector
		".main", ".main-content", ".content", ".primary", ".article", ".story", ".post", ".news-content",
		".article-content", ".story-content", ".content-body", ".main-article", ".primary-content",
		".page-content", ".site-content", ".entry-content", ".post-content", ".article-body", ".story-body",
		".news-body", ".content-wrap", ".content-area", ".main-column", ".primary-column",
		".container .content", ".wrapper .main", ".page-wrapper .content", ".container main",

		// standby selector
		"div[class*='content']:not([class*='nav']):not([class*='sidebar']):not([class*='footer'])",
		"div[id*='content']:not([id*='nav']):not([id*='sidebar']):not([id*='footer'])",
		"section[class*='main']:not([class*='nav']):not([class*='sidebar'])",
		"div[class*='article']:not([class*='related']):not([class*='recommended'])",
	}
	NewsTitleSelectors = []string{
		"h1.title", "h1.article-title", "h1.news-title", "h1.post-title",
		".article-title h1", ".news-title h1", ".post-title h1",
		"h1", "h2", "h3", "h4", "h5", "h6",
		".title", ".headline", ".news-title", ".article-headline", ".article-title",
		".news-headline", ".post-headline", ".story-title", ".post-title", ".story-title", ".content-title",
		".item-title", ".list-title", ".link-title",
		"[class*='title']", "[class*='headline']", "[class*='subject']", "[id*='title']", "[id*='headline']",
		"meta[property='og:title']", "meta[name='title']",
		"a[title]", "a", "div",
	}
	NewsSummarySelectors = []string{
		".summary", ".excerpt", ".description", ".abstract",
		".intro", ".lead", ".subtitle", ".sub-title",
		".news-summary", ".article-summary", ".post-summary",
		".story-summary", ".content-summary", ".item-summary",
		".brief", ".digest", ".outline", ".preview",
		"[class*='summary']", "[class*='excerpt']", "[class*='description']",
		"[class*='abstract']", "[class*='intro']", "[class*='brief']",
		"p", "div",
	}
	NewsContentSelectors = []string{
		// id selector
		"#content p", "#story-body p", "#article-body p", "#content-body p", "#entry-content p", "#post-body p",
		"#article-text p", "#news-body p", "#main-text p",

		// special case
		"div.paragraph", "div.text-block", "div.story-paragraph", "div.article-paragraph", "span.paragraph",

		// class selector
		"article p:not(.caption):not(.meta):not(.byline):not(.share):not(.related)",
		".content p:not(.sidebar-text):not(.widget-text):not(.advertisement):not(.nav-text)",
		".story p:not(.story-tools):not(.story-meta):not(.story-share):not(.story-related)",

		"article .entry-content p", "article .post-content p", "article .article-body p", "main .content-body p",
		"section .story-body p", ".content p", ".story-body p", ".article-body p", ".content-body p", ".entry-content p",
		".post-content p", ".news-body p", ".article-content p", ".story-content p", ".main-text p", ".body-text p",

		"p:not(.caption):not(.meta):not(.byline):not(.timestamp)", "p:not(.related):not(.recommended):not(.advertisement)",
		"p:not(.share):not(.social):not(.tags):not(.category)", "p:not(.navigation):not(.breadcrumb):not(.footer-text)",
		"p:not(.sidebar-text):not(.widget-text):not(.comment)", ".story-body p:not(.caption):not(.meta)",
		".article-body p:not(.byline):not(.share)", ".content-body p:not(.related):not(.ads)",
		".entry-content p:not(.tags):not(.category)",

		"div[class*='content'] p:not([class*='meta']):not([class*='share'])",
		"div[class*='body'] p:not([class*='caption']):not([class*='byline'])",
		"div[class*='text'] p:not([class*='related']):not([class*='ads'])",
		"section[class*='story'] p:not([class*='tools']):not([class*='nav'])",

		// structured selector
		"p", "div p", "article p", "section p", "main p",
	}
	NewsImageSelectors = []string{
		".article-content img", ".news-content img", ".post-content img", ".content img",
		".article-body img", ".news-body img", ".story-body img", ".entry-content img",
		".main-content img", "article img", "main img", ".main img",
		"img[src]", "img", "picture[srcset]", "picture source",
	}
	NewsAuthorSelectors = []string{
		".author", ".writer", ".reporter", ".by-author", ".post-author",
		".article-author", ".news-author", ".story-author",
		".byline", ".by-line", ".author-name", ".writer-name",
		"[class*='author']", "[class*='writer']", "[class*='reporter']",
		"[class*='byline']", "[id*='author']", "[id*='writer']",
		"meta[name='author']", "meta[property='author']",
	}
	NewsTimeSelectors = []string{
		// structured data selector
		"time[datetime]", "[datetime]", "time[data-timestamp]", "[data-timestamp]", "time[data-time]", "[data-time]",

		// semantic selector
		"time", ".publish-time", ".publication-time", ".post-time", ".article-time",
		".news-time", ".date-time", ".timestamp", ".publish-date",
		".post-date", ".article-date", ".news-date", ".date", ".time",

		// class selector
		"[class*='time']", "[class*='date']", "[class*='publish']", "[class*='pub']",
		"[class*='post']", "[id*='time']", "[id*='date']", "[id*='publish']", "[id*='pub']",

		// meta label
		"meta[name='article:published_time']", "meta[property='article:published_time']",
		"meta[name='pubdate']", "meta[property='pubdate']",
		"meta[name='date']", "meta[property='date']",
	}
)

// attributes list
var (
	ImageAttrs     = []string{Attr_src, Attr_srcset, Attr_data_src, Attr_data_original}
	TimeAttributes = []string{
		"datetime", "data-timestamp", "data-time", "data-date", "timestamp", "content", "value", "title",
	}
)

// validate news indicator
var (
	InvalidNewsLinkPrefixes = []string{"javascript:", "mailto:", "tel:", "#", "void(0)"}
	ExcludeNewsLinkPatterns = []string{
		`\.(jpg|jpeg|png|gif|bmp|webp|svg)$`,
		`\.(pdf|doc|docx|xls|xlsx|ppt|pptx)$`,
		`\.(mp3|mp4|avi|mov|wmv|flv)$`,
		`\.(zip|rar|7z|tar|gz)$`,
		`/(tag|tags|category|categories|author|search|login|register|about|contact|help|privacy|terms)/`,
		`\?(tag|tags|category|categories|author|search)=`,
	}
	ImageLinkPatterns = []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg",
		"/image", "/img", "/pic", "/photo", "/upload",
	}
)

// IsNewsLink checks if the given href is a valid news link.
func IsNewsLink(href string) bool {
	if href == "" {
		return false
	}

	if slices.ContainsFunc(InvalidNewsLinkPrefixes, func(prefix string) bool { return strings.HasPrefix(href, prefix) }) {
		return false
	}

	for _, pattern := range ExcludeNewsLinkPatterns {
		if matched, _ := regexp.MatchString(pattern, href); matched {
			return false
		}
	}

	return true
}

// IsNewsImageLink checks if the given href is a valid news image link.
func IsNewsImageLink(href string) bool {
	if href == "" {
		return false
	}

	href = strings.ToLower(href)

	for _, pattern := range ImageLinkPatterns {
		if strings.Contains(href, pattern) {
			return true
		}
	}

	return false
}
