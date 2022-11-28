package middleware

import (
	"github.com/farseer-go/webapi/context"
	"strings"
)

type Routing struct {
	IMiddleware
}

func (receiver *Routing) Invoke(httpContext *context.HttpContext) {
	// 检查method
	if strings.ToUpper(httpContext.HttpRoute.Method) != httpContext.Method {
		// 响应码
		httpContext.HttpResponse.WriteCode(405)
		return
	}
	receiver.IMiddleware.Invoke(httpContext)
}
