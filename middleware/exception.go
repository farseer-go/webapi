package middleware

import (
	"fmt"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi/context"
)

// exception 异常中间件（默认加载）
type exceptionMiddleware struct {
	IMiddleware
}

func (receiver *exceptionMiddleware) Invoke(httpContext *context.HttpContext) {
	exception.Try(func() {
		receiver.IMiddleware.Invoke(httpContext)
	}).CatchWebException(func(exp *exception.WebException) {
		// 响应码
		httpContext.Response.StatusCode = exp.StatusCode
		httpContext.Response.BodyString = exp.Message
		httpContext.Response.BodyBytes = []byte(exp.Message)
		httpContext.Exception = exp.Message
	}).CatchException(func(exp any) {
		// 响应码
		httpContext.Response.StatusCode = 500
		httpContext.Response.BodyString = fmt.Sprint(exp)
		httpContext.Response.BodyBytes = []byte(httpContext.Response.BodyString)
		httpContext.Exception = exp
		flog.Warningf("[%s]%s 发生错误：%s", httpContext.Method, httpContext.URI.Url, httpContext.Response.BodyString)
	})
}
