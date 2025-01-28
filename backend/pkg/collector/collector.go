package collector

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// NewCollector creates a new collector instance.
func NewCollector() *colly.Collector {
	c := colly.NewCollector(
	// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
		// Parallelism: 2,
	})

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	extensions.RandomMobileUserAgent(c)
	extensions.Referer(c)

	return c
}
