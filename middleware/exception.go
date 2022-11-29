package middleware

import (
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
)

// Exception 异常中间件（默认加载）
type Exception struct {
	IMiddleware
}

func (receiver *Exception) Invoke(httpContext *context.HttpContext) {
	exception.Try(func() {
		receiver.IMiddleware.Invoke(httpContext)
		// 响应码
		if httpContext.HttpResponse.StatusCode == 0 {
			httpContext.HttpResponse.StatusCode = 200
		}
	}).CatchWebException(func(exp *exception.WebException) {
		// 响应码
		httpContext.HttpResponse.StatusCode = exp.StatusCode
		httpContext.HttpResponse.BodyString = exp.Message
		httpContext.HttpResponse.BodyBytes = []byte(exp.Message)
		httpContext.Exception = exp.Message
	}).CatchException(func(exp any) {
		// 响应码
		httpContext.HttpResponse.StatusCode = 500
		httpContext.HttpResponse.BodyString = exp.(string)
		httpContext.HttpResponse.BodyBytes = []byte(httpContext.HttpResponse.BodyString)
		httpContext.Exception = exp
	})
}