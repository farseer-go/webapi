package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
)

type ApiResponse struct {
	context.IMiddleware
}

func (receiver *ApiResponse) Invoke(httpContext *context.HttpContext) {
	// ActionResult类型，不做ApiResponse解析
	if httpContext.IsActionResult() {
		receiver.IMiddleware.Invoke(httpContext)
		return
	}

	var apiResponse core.ApiResponse[any]
	catch := exception.Try(func() {
		receiver.IMiddleware.Invoke(httpContext)

		var returnVal any
		// 只有一个返回值
		bodyLength := len(httpContext.Response.Body)
		if bodyLength == 1 {
			returnVal = httpContext.Response.Body[0]
		} else if bodyLength > 1 {
			// 多个返回值，则转成数组Json
			lst := collections.NewListAny()
			for i := 0; i < bodyLength; i++ {
				lst.Add(httpContext.Response.Body[i])
			}
			returnVal = lst
		}
		statusCode, statusMessage := httpContext.Response.GetStatus()
		apiResponse = core.Success[any](statusMessage, returnVal)
		apiResponse.StatusCode = statusCode
		apiResponse.Status = statusCode == 200
	})

	catch.CatchWebException(func(exp exception.WebException) {
		// 响应码
		apiResponse = core.Error[any](exp.Message, exp.StatusCode)
	})

	catch.CatchException(func(exp any) {
		// 响应码
		apiResponse = core.Error[any](fmt.Sprint(exp), http.StatusInternalServerError)
	})

	traceContext := httpContext.Data.Get("Trace").(*trace.TraceContext)
	apiResponse.TraceId = traceContext.TraceId
	apiResponse.ElapsedMilliseconds = (time.Now().UnixMicro() - traceContext.StartTs) / 1000
	httpContext.Route.IsGoBasicType = false
	httpContext.Response.Body = []any{apiResponse}
}
