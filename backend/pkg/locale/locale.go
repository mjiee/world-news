package locale

import (
	"embed"
	"encoding/json"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	En = "en"
	Zh = "zh"
)

var (
	once   sync.Once
	bundle *i18n.Bundle
)

//go:embed locales/locale.*.json
var LocaleFS embed.FS

// LoadLocaleFile loads the locale file
func LoadLocaleFile() (*i18n.Bundle, error) {
	var err error

	once.Do(func() {
		bundle = i18n.NewBundle(language.English)

		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		_, err = bundle.LoadMessageFileFS(LocaleFS, "locales/locale.en.json")
		if err != nil {
			return
		}

		_, err = bundle.LoadMessageFileFS(LocaleFS, "locales/locale.zh.json")
		if err != nil {
			return
		}
	})

	return bundle, err
}

// Localize localize
func Localize(localizer *i18n.Localizer, messageId string, templateData ...map[string]string) string {
	config := &i18n.LocalizeConfig{MessageID: messageId}

	if len(templateData) > 0 {
		config.TemplateData = templateData[0]
	}

	msg, err := localizer.Localize(config)
	if err != nil {
		return messageId
	}

	return msg
}
