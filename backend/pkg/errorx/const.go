package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "Internal Server Error.")
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
	CrawlingRecordNotFound = NewBasicError(103011, "Crawling record not found.")
)
