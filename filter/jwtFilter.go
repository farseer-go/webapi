package filter

import (
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
	"golang.org/x/net/websocket"
)

type JwtFilter struct {
}

func (receiver JwtFilter) OnActionExecuting(httpContext *context.HttpContext) {
	if !httpContext.Jwt.Valid() {
		// ws协议验证失败时，先主动发一条消息，然后立即关闭
		if httpContext.WebsocketConn != nil {
			_ = websocket.JSON.Send(httpContext.WebsocketConn, core.ApiResponseStringError(context.InvalidMessage, context.InvalidStatusCode))
		}
		exception.ThrowWebExceptionf(context.InvalidStatusCode, context.InvalidMessage)
	}
}

func (receiver JwtFilter) OnActionExecuted(httpContext *context.HttpContext) {
}
