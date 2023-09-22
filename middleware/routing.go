package middleware

import (
	"github.com/farseer-go/webapi/context"
)

type routing struct {
	context.IMiddleware
}

func (receiver *routing) Invoke(httpContext *context.HttpContext) {
	// 检查method
	if httpContext.Method != "OPTIONS" && !httpContext.Route.Method.Contains(httpContext.Method) {
		// 响应码
		httpContext.Response.Error405()
		return
	}
	receiver.IMiddleware.Invoke(httpContext)
}
