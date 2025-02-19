//go:build web

package httpx

import (
	"net/http"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/locale"

	"github.com/gin-gonic/gin"
)

// WebResp is a function that handles the response of web application.
func WebResp(c *gin.Context, result any, err error) {
	var resp = Ok(result)

	if err != nil {
		basicErr, ok := err.(*errorx.BasicError)
		if ok {
			if basicErr.GetErr() != nil {
				c.Error(basicErr.GetErr())
			}

			resp = NewResponse(basicErr.GetCode(), locale.WebLocalize(c, basicErr.GetMessage()))
		} else if !ok {
			c.Error(err)

			resp = NewResponse(errorx.InternalError.GetCode(),
				locale.WebLocalize(c, errorx.InternalError.GetMessage()))
		}
	}

	c.JSON(http.StatusOK, resp)
}
