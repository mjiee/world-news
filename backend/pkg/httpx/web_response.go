//go:build web

package httpx

import (
	"net/http"

	"github.com/mjiee/world-news/backend/pkg/errorx"

	"github.com/gin-gonic/gin"
)

// WebResp is a function that handles the response of web application.
func WebResp(c *gin.Context, result any, err error) {
	resp := Resp(result, err)

	if err != nil {
		basicErr, ok := err.(*errorx.BasicError)
		if ok && basicErr.GetErr() != nil {
			c.Error(basicErr.GetErr())
		} else if !ok {
			c.Error(err)
		}
	}

	c.JSON(http.StatusOK, resp)
}
