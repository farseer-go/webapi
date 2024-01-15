package middleware

import (
	"fmt"
	"github.com/farseer-go/fs/exception"
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
		httpContext.Response.SetHttpCode(exp.StatusCode)
	}).CatchException(func(exp any) {
		switch e := exp.(type) {
		case error:
			httpContext.Exception = e
		default:
			httpContext.Exception = fmt.Errorf("%s", e)
		}
		// 响应码
		httpContext.Response.Write([]byte(fmt.Sprint(exp)))
		httpContext.Response.SetHttpCode(http.StatusInternalServerError)
	})
}
