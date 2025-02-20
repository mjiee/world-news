package valueobject

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// SystemConfigKey system config key
type SystemConfigKey string

// system config key
const (
	NewsWebsiteCollectionKey SystemConfigKey = "newsWebsiteCollection" // news website collection
	NewsWebsiteKey           SystemConfigKey = "newsWebsite"           // news website
	NewsTopicKey             SystemConfigKey = "newsTopic"             // news topic
	LanguageKey              SystemConfigKey = "language"              // language
)

func (s SystemConfigKey) String() string {
	return string(s)
}

// MarshalValue marshal config value
func (s SystemConfigKey) MarshalValue(value any) (string, error) {
	switch s {
	case LanguageKey:
		return value.(string), nil
	}

	data, err := json.Marshal(value)

	return string(data), errors.WithStack(err)
}

// UnmarshalValue unmarshal config value
func (s SystemConfigKey) UnmarshalValue(value string) (any, error) {
	switch s {
	case NewsWebsiteCollectionKey, NewsWebsiteKey:
		var (
			data []*NewsWebsite
			err  = json.Unmarshal([]byte(value), &data)
		)

		return data, errors.WithStack(err)
	case NewsTopicKey:
		var (
			data []string
			err  = json.Unmarshal([]byte(value), &data)
		)

		return data, errors.WithStack(err)
	case LanguageKey:
		return value, nil
	}

	var (
		data any
		err  = json.Unmarshal([]byte(value), &data)
	)

	return data, errors.WithStack(err)
}
