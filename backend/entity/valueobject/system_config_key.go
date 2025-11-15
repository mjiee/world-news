package valueobject

// SystemConfigKey system config key
type SystemConfigKey string

// system config key
const (
	NewsWebsiteCollectionKey SystemConfigKey = "newsWebsiteCollections" // news website collection
	NewsWebsiteKey           SystemConfigKey = "newsWebsites"           // news website
	NewsTopicKey             SystemConfigKey = "newsTopics"             // news topic
	LanguageKey              SystemConfigKey = "language"               // language
	RemoteService            SystemConfigKey = "remoteService"          // remote service
	TranslaterKey            SystemConfigKey = "translater"             // translater
	InvalidNewsWebsiteKey    SystemConfigKey = "invalidNewsWebsite"     // invalid news website
	TextToSpeechAIKey        SystemConfigKey = "textToSpeechAI"         // text to speech ai
	TextAIKey                SystemConfigKey = "textAI"                 // openai
	NewsCritiquePromptKey    SystemConfigKey = "newsCritiquePrompt"     // news critique prompt
	PodcastScriptPromptKey   SystemConfigKey = "podcastScriptPrompt"    // podcast script prompt
)

func (s SystemConfigKey) String() string {
	return string(s)
}
