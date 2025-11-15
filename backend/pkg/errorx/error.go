package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "error.internalError")
	ParamsError   = NewBasicError(http.StatusBadRequest, "error.paramsError")
	Unauthorized  = NewBasicError(http.StatusUnauthorized, "error.unauthorized")
)

// system config error
var (
	SystemConfigNotFound      = NewBasicError(101011, "error.systemConfigNotFound")
	OpenaiConfigNotFound      = NewBasicError(101012, "error.openaiConfigNotFound")
	TtsAiConfigNotFound       = NewBasicError(101013, "error.ttsAiConfigNotFound")
	TranslaterConfigNotFound  = NewBasicError(101014, "error.translaterConfigNotFound")
	CritiquePromptNotFound    = NewBasicError(101015, "error.critiquePromptNotFound")
	NewsWebsiteConfigNotFound = NewBasicError(101016, "error.newsWebsiteConfigNotFound")
	PodcastPromptNotFound     = NewBasicError(101017, "error.podcastPromptNotFound")
	PodcastVoiceNotFound      = NewBasicError(101018, "error.podcastVoiceNotFound")
)

// news error
var (
	NewsNotFound = NewBasicError(102011, "error.newsNotFound")
)

// crawling error
var (
	CrawlingRecordNotFound       = NewBasicError(103011, "error.crawlingRecordNotFound")
	HasProcessingTasks           = NewBasicError(103012, "error.hasProcessingTasks")
	UpdateRecordStatusNotAllowed = NewBasicError(103013, "error.updateRecordStatusNotAllowed")
)

// podcast error
var (
	PodcastTaskNotFound     = NewBasicError(104011, "error.podcastTaskNotFound")
	PodcastGenerationFailed = NewBasicError(104012, "error.podcastGenerationFailed")
	PodcastScriptNotFound   = NewBasicError(104013, "error.podcastScriptNotFound")
)
