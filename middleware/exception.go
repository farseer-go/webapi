package middleware

import (
	"net/http"

	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
	"golang.org/x/net/websocket"
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
		// ws协议先主动发一条消息，然后立即关闭
		if httpContext.WebsocketConn != nil {
			_ = websocket.JSON.Send(httpContext.WebsocketConn, core.ApiResponseStringError(exp.Message, exp.StatusCode))
		}
		// 响应码
		httpContext.Response.Write([]byte(exp.Message))
		httpContext.Response.SetHttpCode(exp.StatusCode)
	}).CatchException(func(exp any) {
		// ws协议先主动发一条消息，然后立即关闭
		if httpContext.WebsocketConn != nil {
			_ = websocket.JSON.Send(httpContext.WebsocketConn, core.ApiResponseStringError(httpContext.Exception.Error(), http.StatusInternalServerError))
		}

		// 响应码
		httpContext.Response.Write([]byte(httpContext.Exception.Error()))
		httpContext.Response.SetHttpCode(http.StatusInternalServerError)
	})
}
