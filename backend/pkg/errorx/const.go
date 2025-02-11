package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "internal server error")
)

// system config error
var (
	SystemConfigNotFound = NewBasicError(101011, "system config not found")
)

// news error
var (
	NewsNotFound = NewBasicError(102011, "news not found")
)

// crawling error
var (
	CrawlingRecordNotFound = NewBasicError(103011, "crawling record not found")
)
