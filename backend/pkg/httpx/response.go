package httpx

// Response is a public response struct.
type Response struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

// NewResponse creates and returns a new Response object
func NewResponse(code uint32, message string, result any) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Result:  result,
	}
}

// OK is a convenience function used to create a response object with a successful status (code 0 and no message).
func OK(result any) *Response {
	return NewResponse(0, "", result)
}
