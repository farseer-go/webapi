package middleware

import "github.com/farseer-go/webapi/context"

type Session struct {
	context.IMiddleware
}

func (receiver *Session) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)
}
