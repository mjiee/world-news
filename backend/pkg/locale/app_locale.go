package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var appLocalizer *i18n.Localizer

// SetAppLocalizer sets the app localizer
func SetAppLocalizer(lang string) error {
	if lang == "" {
		lang = En
	}

	bundle, err := LoadLocaleFile()
	if err != nil {
		return err
	}

	appLocalizer = i18n.NewLocalizer(bundle, lang)

	return nil
}

// AppLocalize app localize
func AppLocalize(messageId string, templateData ...map[string]string) string {
	if appLocalizer == nil {
		SetAppLocalizer(En)
	}

	return Localize(appLocalizer, messageId, templateData...)
}
