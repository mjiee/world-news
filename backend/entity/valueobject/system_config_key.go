package valueobject

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
