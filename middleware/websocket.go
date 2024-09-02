package middleware

import (
	"github.com/farseer-go/webapi/context"
	"net/http"
)

type Websocket struct {
	context.IMiddleware
}

func (receiver *Websocket) Invoke(httpContext *context.HttpContext) {
	httpContext.Response.SetHttpCode(http.StatusOK)
	httpContext.Response.SetStatusCode(http.StatusOK)

	// 下一步：exceptionMiddleware
	receiver.IMiddleware.Invoke(httpContext)

	_ = httpContext.WebsocketConn.WriteClose(httpContext.Response.GetHttpCode())
}
