package valueobject

// SystemConfigKey system config key
type SystemConfigKey string

// system config key
const (
	NewsWebsiteCollectionKey SystemConfigKey = "newsWebsiteCollections" // news website collection
	NewsWebsiteKey           SystemConfigKey = "newsWebsites"           // news website
	NewsTopicKey             SystemConfigKey = "newsTopics"             // news topic
	LanguageKey              SystemConfigKey = "language"               // language
	OpenAIKey                SystemConfigKey = "openAI"                 // openai
	TranslaterKey            SystemConfigKey = "translater"             // translater
	InvalidNewsWebsiteKey    SystemConfigKey = "invalidNewsWebsite"     // invalid news website
)

func (s SystemConfigKey) String() string {
	return string(s)
}
