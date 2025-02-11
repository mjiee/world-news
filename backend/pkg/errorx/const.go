package errorx

import "net/http"

// basic error
var (
	InternalError = NewBasicError(http.StatusInternalServerError, "internal server error")
)

// system config error
