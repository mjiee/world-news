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
	SystemConfigNotFound = NewBasicError(101011, "error.systemConfigNotFound")
)

// news error
var (
	NewsNotFound             = NewBasicError(102011, "error.newsNotFound")
	OpenaiConfigNotFound     = NewBasicError(102012, "error.openaiConfigNotFound")
	TranslaterConfigNotFound = NewBasicError(102013, "error.translaterConfigNotFound")
)

// crawling error
var (
	CrawlingRecordNotFound       = NewBasicError(103011, "error.crawlingRecordNotFound")
	HasProcessingTasks           = NewBasicError(103012, "error.hasProcessingTasks")
	NewsWebsiteConfigNotFound    = NewBasicError(103013, "error.newsWebsiteConfigNotFound")
	NewsTopicConfigNotFound      = NewBasicError(103014, "error.newsTopicConfigNotFound")
	UpdateRecordStatusNotAllowed = NewBasicError(103015, "error.updateRecordStatusNotAllowed")
)
