package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "error.InternalError")
	ParamsError   = NewBasicError(http.StatusBadRequest, "error.ParamsError")
)

// system config error
var (
	SystemConfigNotFound = NewBasicError(101011, "error.ystemConfigNotFound")
)

// news error
var (
	NewsNotFound = NewBasicError(102011, "error.NewsNotFound")
)

// crawling error
var (
	CrawlingRecordNotFound    = NewBasicError(103011, "error.CrawlingRecordNotFound")
	HasProcessingTasks        = NewBasicError(103012, "error.HasProcessingTasks")
	NewsWebsiteConfigNotFound = NewBasicError(103013, "error.NewsWebsiteConfigNotFound")
	NewsTopicConfigNotFound   = NewBasicError(103014, "error.NewsTopicConfigNotFound")
)
