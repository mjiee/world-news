package valueobject

import (
	"regexp"
	"slices"
	"strings"
)

// common selectors
var (
	ExcludeSelectors = []string{
		"header", ".header", "#header", ".footer", ".footer", "#footer", "aside", ".aside", ".sidebar", ".side-bar",
		"nav", ".nav", "navigation", ".navigation", ".menu", ".top-menu", ".main-menu",
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
		"[class*='sidebar']", "[class*='aside']", "[class*='footer']",
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
		"p", ".summary", ".excerpt", ".description", ".abstract",
		".intro", ".lead", ".subtitle", ".sub-title",
		".news-summary", ".article-summary", ".post-summary",
		".story-summary", ".content-summary", ".item-summary",
		".brief", ".digest", ".outline", ".preview",
		"[class*='summary']", "[class*='excerpt']", "[class*='description']",
		"[class*='abstract']", "[class*='intro']", "[class*='brief']",
		"div",
	}
	NewsContentSelectors = []string{
		"p", ".article-content", ".news-content", ".post-content", ".story-content",
		".content", ".article-body", ".news-body", ".post-body",
		".story-body", ".entry-content", ".main-content", ".text-content",
		".article-text", ".news-text", ".post-text", ".story-text",
		".detail-content", ".full-content", ".rich-content",
		"[class*='content']", "[class*='body']", "[class*='text']",
		"[class*='article']", "[class*='news']", "[class*='post']",
		"[id*='content']", "[id*='body']", "[id*='text']",
		"div",
	}
	NewsImageSelectors = []string{
		".article-content img", ".news-content img", ".post-content img", ".content img",
		".article-body img", ".news-body img", ".story-body img", ".entry-content img",
		".main-content img", "article img", "main img", ".main img",
		"img[src]", "picture[srcset]", "picture source",
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
