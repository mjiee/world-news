//go:build web

package tracex

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceparent = "Traceparent"
)

// WebTracer is a middleware that traces the request and response.
func WebTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if !trace.SpanFromContext(ctx).SpanContext().IsValid() {
			ctx = ExtractTraceFromRequest(ctx, c.Request)

			c.Request = c.Request.WithContext(ctx)
		}

		c.Writer.Header().Set(traceparent, ExtractTraceparent(ctx))

		c.Next()
	}
}
