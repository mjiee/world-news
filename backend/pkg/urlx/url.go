package urlx

import (
	"fmt"
	"net/url"
	"strings"
)

// UrlPrefixHandle url prefix handle
func UrlPrefixHandle(webUrl string, reqUrl *url.URL) string {
	if strings.HasPrefix(webUrl, "http") {
		return webUrl
	}

	if strings.HasPrefix(webUrl, "//") {
		return fmt.Sprintf("%s:%s", reqUrl.Scheme, webUrl)
	}

	webUrl = strings.TrimLeft(webUrl, "/")

	if !strings.Contains(webUrl, reqUrl.Host) {
		webUrl = copyUrl(reqUrl).JoinPath(webUrl).String()
	}

	if !strings.HasPrefix(webUrl, "http") {
		webUrl = fmt.Sprintf("%s://%s", reqUrl.Scheme, webUrl)
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
