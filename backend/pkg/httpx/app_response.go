package httpx

import (
	"context"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/logx"
	"github.com/mjiee/world-news/backend/pkg/tracex"
)

// AppRespHandle is a function that handles the response of desktop app.
func AppRespHandle(ctx context.Context, path string, req, result any, err error) *Response {
	var (
		resp = Resp(result, err)
		data = &logx.LogData{
			Request:  req,
			Response: resp,
			Duration: tracex.CalculateDuration(ctx),
		}
	)

	if err != nil {
		basicErr, ok := err.(*errorx.BasicError)
		if ok && basicErr.GetErr() != nil {
			logx.WithContext(ctx).Error(path, basicErr.GetErr())
		} else if !ok {
			logx.WithContext(ctx).Error(path, err)
		}
	}

	logx.WithContext(ctx).Info(path, data)

	return resp
}
