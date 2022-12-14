package middleware

import (
	"github.com/farseer-go/webapi/context"
	"strings"
)

type routing struct {
	IMiddleware
}

func (receiver *routing) Invoke(httpContext *context.HttpContext) {
	// 检查method
	if httpContext.Method != "OPTIONS" && strings.ToUpper(httpContext.Route.Method) != httpContext.Method {
		// 响应码
		httpContext.Response.StatusCode = 405
		return
	}
	receiver.IMiddleware.Invoke(httpContext)
}
