package urlx

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/mjiee/world-news/backend/pkg/textx"
)

// ExtractHostFromURL extract host from url
func ExtractHostFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	if u.Host == "" {
		return urlStr
	}

	return u.Hostname()
}

// ExtractSecondLevelDomain extract second level domain from url
func ExtractSecondLevelDomain(urlStr string) string {
	host := ExtractHostFromURL(urlStr)

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return host
	}

	parts = parts[:len(parts)-1]

	if len(parts) < 2 {
		return strings.Join(parts, ".")
	}

	return strings.Join(parts[1:], ".")
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

	if strings.HasPrefix(href, "http") {
		return href
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	if strings.HasPrefix(href, ":") {
		return fmt.Sprintf("%s%s", base.Scheme, href)
	}

	if strings.HasPrefix(href, "//") {
		return fmt.Sprintf("%s:%s", base.Scheme, href)
	}

	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, href)
	}

	paths := strings.Split(href, "/")

	for _, prefix := range buildHostPrefix(base.Host, "com") {
		if strings.Contains(paths[0], prefix) || strings.Contains(prefix, paths[0]) {
			return fmt.Sprintf("%s://%s", base.Scheme, href)
		}
	}

	return fmt.Sprintf("%s://%s/%s", base.Scheme, base.Host, href)
}

// buildHostPrefix build host prefix
func buildHostPrefix(hosts ...string) []string {
	prefixs := make([]string, 0, len(hosts))

	for _, host := range hosts {
		prefixs = append(prefixs, host)

		if strings.HasPrefix(host, "www.") && len(host) > 4 {
			host = host[4:]

			prefixs = append(prefixs, host)
		}

		parts := strings.Split(host, ".")
		if len(parts) < 2 {
			continue
		}

		prefixs = append(prefixs, strings.Join(parts[:len(parts)-1], "."))
	}

	return prefixs
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
