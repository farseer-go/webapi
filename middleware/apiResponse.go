package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
)

type ApiResponse struct {
	IMiddleware
}

func (receiver *ApiResponse) Invoke(httpContext *context.HttpContext) {

	var apiResponse core.ApiResponse[any]
	exception.Try(func() {
		receiver.IMiddleware.Invoke(httpContext)

		var returnVal any
		// 只有一个返回值
		if len(httpContext.HttpResponse.Body) == 1 {
			returnVal = httpContext.HttpResponse.Body[0].Interface()
		} else {
			// 多个返回值，则转成数组Json
			lst := collections.NewListAny()
			for i := 0; i < len(httpContext.HttpResponse.Body); i++ {
				lst.Add(httpContext.HttpResponse.Body[i].Interface())
			}
			returnVal = lst
		}
		apiResponse = core.Success[any]("成功", returnVal)
	}).CatchWebException(func(exp *exception.WebException) {
		// 响应码
		httpContext.HttpResponse.StatusCode = exp.StatusCode
		httpContext.Exception = exp.Message
		apiResponse = core.Error[any](exp.Message, httpContext.HttpResponse.StatusCode)
	}).CatchException(func(exp any) {
		// 响应码
		httpContext.HttpResponse.StatusCode = 500
		httpContext.Exception = exp
		apiResponse = core.Error[any](exp.(string), httpContext.HttpResponse.StatusCode)
	})

	httpContext.HttpResponse.BodyBytes = apiResponse.ToBytes()
	httpContext.HttpResponse.BodyString = string(httpContext.HttpResponse.BodyBytes)
}
