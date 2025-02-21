package urlx

import (
	"net/url"
	"strings"
)

// UrlPrefixHandle url prefix handle
func UrlPrefixHandle(webUrl string, reqUrl *url.URL) string {
	if strings.HasPrefix(webUrl, "http") {
		return webUrl
	}

	if !strings.Contains(webUrl, reqUrl.Host) {
		webUrl = copyUrl(reqUrl).JoinPath(webUrl).String()
	}

	if strings.HasPrefix(webUrl, "/") {
		return ""
	}

	return webUrl
}

// copyUrl copy url
func copyUrl(src *url.URL) *url.URL {
	return &url.URL{
		Host:   src.Host,
		Scheme: src.Scheme,
	}
}
