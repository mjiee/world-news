package httpx

import "github.com/mjiee/world-news/backend/pkg/errorx"

// Response is a public response struct.
type Response struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

// NewResponse creates and returns a new Response object
func NewResponse(code uint32, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

// Fail is a convenience function used to create a response object with a failed status (code 1 and no message).
func Fail(err error) *Response {
	basicErr, ok := err.(*errorx.BasicError)
	if ok {
		return NewResponse(basicErr.GetCode(), basicErr.GetMessage())
	}

	return NewResponse(errorx.InternalError.GetCode(), errorx.InternalError.GetMessage())
}
