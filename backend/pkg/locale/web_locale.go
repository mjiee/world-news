//go:build web

package locale

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/mjiee/world-news/backend/pkg/logx"
)

const (
	acceptLanguage  = "Accept-Language"
	webLocalizerKey = "localizer"
)

// WebLocale is a middleware that sets the localizer for the request.
func WebLocale() gin.HandlerFunc {
	bundle, err := LoadLocaleFile()
	if err != nil {
		logx.Fatal("load locale file failed", err)
	}

	return func(c *gin.Context) {
		lang := c.GetHeader(acceptLanguage)

		if lang == "" {
			lang = c.DefaultQuery("lang", En)
		}

		c.Set(webLocalizerKey, i18n.NewLocalizer(bundle, lang))

		c.Next()
	}
}

// WebLocalize web localize
func WebLocalize(c *gin.Context, messageId string, templateData ...map[string]string) string {
	localizer, ok := c.MustGet(webLocalizerKey).(*i18n.Localizer)
	if !ok {
		return messageId
	}

	return Localize(localizer, messageId, templateData...)
}
