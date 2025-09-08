//go:build web

package httpx

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ParseRequest parse request
func ParseRequest[T any](c *gin.Context) (context.Context, T, error) {
	var (
		ctx = c.Request.Context()
		req T
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		return ctx, req, err
	}

	return ctx, req, nil
}
