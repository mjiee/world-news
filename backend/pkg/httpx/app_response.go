//go:build !web

package httpx

import (
	"context"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/locale"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
)

// AppResp is a function that handles the response of desktop app.
func AppResp(ctx context.Context, path string, req, result any, err error) *Response {
	var (
		resp = Ok(result)
		data = &logx.LogData{
			Request:  req,
			Response: resp,
			Duration: tracex.CalculateDuration(ctx),
		}
	)

	if err != nil {
		basicErr, ok := err.(*errorx.BasicError)
		if ok {
			if basicErr.GetErr() != nil {
				logx.WithContext(ctx).Error(path, basicErr.GetErr())
			}

			resp = NewResponse(basicErr.GetCode(), locale.AppLocalize(basicErr.GetMessage()))
		} else {
			logx.WithContext(ctx).Error(path, err)

			resp = NewResponse(errorx.InternalError.GetCode(), locale.AppLocalize(errorx.InternalError.GetMessage()))
		}
	}

	logx.WithContext(ctx).Info(path, data)

	if resp == nil {
		resp = Ok(nil)
	}

	return resp
}
