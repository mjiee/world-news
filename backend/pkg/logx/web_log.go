//go:build web

package logx

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	contentType = "Content-Type"
)

// WebLogger is a middleware that logs the request and response.
func WebLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				WithContext(c.Request.Context()).Error(c.Request.URL.Path, errors.Errorf("%v", err))
			}

			c.AbortWithStatus(http.StatusInternalServerError)
		}()

		var (
			startTime   = time.Now()
			requestBody []byte
			logData     = &LogData{
				Method: c.Request.Method,
			}
		)

		// request body
		if isJsonBody(c.Request.Header) && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)

			logData.Request = string(requestBody)

			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		} else if c.Request.Method == http.MethodGet {
			queryParams := c.Request.URL.Query()

			if len(queryParams) > 0 {
				logData.Request = queryParams.Encode()
			}
		}

		// respose body
		respWriter := &respBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}

		c.Writer = respWriter

		c.Next()

		logFlag := false

		for _, err := range c.Errors {
			WithContext(c.Request.Context()).Error(c.Request.URL.Path, err.Err)
			logFlag = true
		}

		if !logFlag {
			return
		}

		if isJsonBody(c.Writer.Header()) && respWriter.body != nil {
			logData.Response = respWriter.body.String()
		}

		logData.Duration = time.Since(startTime).Milliseconds()

		WithContext(c.Request.Context()).Info(c.Request.URL.Path, logData)
	}
}

// respBodyWriter is a custom ResponseWriter that captures the response body.
type respBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (blw *respBodyWriter) Write(p []byte) (n int, err error) {
	n, err = blw.body.Write(p)
	if err != nil {
		return n, err
	}

	return blw.ResponseWriter.Write(p)
}

// isJsonBody checks if the request/response body is a JSON body.
func isJsonBody(header http.Header) bool {
	return strings.Contains(header.Get(contentType), gin.MIMEJSON)
}
