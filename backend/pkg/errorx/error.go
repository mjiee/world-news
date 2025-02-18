package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "Internal Server Error.")
	ParamsError   = NewBasicError(http.StatusBadRequest, "Request params error.")
)

// system config error
var (
	SystemConfigNotFound = NewBasicError(101011, "System config not found.")
)

// news error
var (
	NewsNotFound = NewBasicError(102011, "News not found.")
)

// crawling error
var (
	CrawlingRecordNotFound    = NewBasicError(103011, "Crawling record not found.")
	HasProcessingTasks        = NewBasicError(103012, "There are still processing tasks. Please try again later.")
	NewsWebsiteConfigNotFound = NewBasicError(103013, "Please configure the news website before crawling.")
	NewsTopicConfigNotFound   = NewBasicError(103014, "Please configure the news topic before crawling.")
)
