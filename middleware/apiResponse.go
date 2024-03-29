package middleware

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"net/http"
	"time"
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
		if len(httpContext.Response.Body) == 1 {
			returnVal = httpContext.Response.Body[0]
		} else {
			// 多个返回值，则转成数组Json
			lst := collections.NewListAny()
			for i := 0; i < len(httpContext.Response.Body); i++ {
				lst.Add(httpContext.Response.Body[i])
			}
			returnVal = lst
		}
		apiResponse = core.Success[any](httpContext.Response.GetStatusMessage(), returnVal)
		apiResponse.StatusCode = httpContext.Response.GetStatusCode()
	})

	catch.CatchWebException(func(exp exception.WebException) {
		// 响应码
		apiResponse = core.Error[any](exp.Message, exp.StatusCode)
	})

	catch.CatchException(func(exp any) {
		switch e := exp.(type) {
		case error:
			httpContext.Exception = e
		default:
			httpContext.Exception = fmt.Errorf("%s", e)
		}
		// 响应码
		apiResponse = core.Error[any](fmt.Sprint(exp), http.StatusInternalServerError)
	})

	traceContext := httpContext.Data.Get("Trace").(trace.ITraceContext)
	apiResponse.TraceId = traceContext.GetTraceId()
	apiResponse.ElapsedMilliseconds = (time.Now().UnixMicro() - traceContext.GetStartTs()) / 1000
	httpContext.Route.IsGoBasicType = false
	httpContext.Response.Body = []any{apiResponse}
}
