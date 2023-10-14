package middleware

import (
	"fmt"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi/context"
	"net/http"
)

// exception 异常中间件（默认加载）
type exceptionMiddleware struct {
	context.IMiddleware
}

func (receiver *exceptionMiddleware) Invoke(httpContext *context.HttpContext) {
	// exceptionMiddleware 与 ApiResponse 中间件是互诉的。
	exception.Try(func() {
		// 下一步：routing
		receiver.IMiddleware.Invoke(httpContext)
	}).CatchWebException(func(exp exception.WebException) {
		// 响应码
		httpContext.Response.Write([]byte(exp.Message))
		httpContext.Exception = exp.Message
		httpContext.Response.SetHttpCode(exp.StatusCode)
	}).CatchException(func(exp any) {
		// 响应码
		httpContext.Response.Write([]byte(fmt.Sprint(exp)))
		httpContext.Exception = exp
		httpContext.Response.SetHttpCode(http.StatusInternalServerError)
		flog.Warningf("[%s]%s 发生错误：%s", httpContext.Method, httpContext.URI.Url, fmt.Sprint(exp))
	})
}
