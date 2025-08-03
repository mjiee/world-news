package urlx

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/mjiee/world-news/backend/pkg/textx"
)

// ExtractDomainFromURL extract domain from url
func ExtractDomainFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	host := u.Host
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	return host
}

// RemoveQueryParams removes the query parameters from a URL string
func RemoveQueryParams(rawURL string) string {
	if idx := strings.Index(rawURL, "?"); idx != -1 {
		return rawURL[:idx]
	}

	return rawURL
}

// NormalizeURL normalize url
func NormalizeURL(baseURL, href string) string {
	href = textx.CleanText(href)

	if href == "" {
		return ""
	}

	// srcset href
	hrefs := strings.Split(href, ",")
	if len(hrefs) > 2 {
		href = hrefs[0]
	}

	hrefs = strings.Fields(textx.CleanText(href))
	if len(hrefs) > 2 {
		href = hrefs[0]
	}

	// replace relative url
	if strings.HasPrefix(href, "/") {
		return resetHref(baseURL, strings.TrimPrefix(href, "/"))
	}

	// replace relative url
	if !strings.HasPrefix(href, "http") {
		return resetHref(baseURL, href)
	}

	return href
}

// resetHref reset href
func resetHref(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return href
	}

	return fmt.Sprintf("%s://%s/%s", baseURL.Scheme, baseURL.Host, strings.TrimPrefix(href, "/"))
}

// invalidNewsLinkPrefixes invalid news link prefixes
var invalidNewsLinkPrefixes = []string{"javascript:", "mailto:", "tel:", "#", "void(0)"}

// IsValidURL checks if the given URL is valid.
func IsValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	if slices.ContainsFunc(invalidNewsLinkPrefixes, func(prefix string) bool { return strings.HasPrefix(rawURL, prefix) }) {
		return false
	}

	if strings.HasPrefix(rawURL, "/") {
		return true
	}

	_, err := url.Parse(rawURL)

	return err == nil
}
