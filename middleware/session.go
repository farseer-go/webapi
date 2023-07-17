package middleware

import "github.com/farseer-go/webapi/context"

type Session struct {
	context.IMiddleware
}

func (receiver *Session) Invoke(httpContext *context.HttpContext) {
	httpContext.Session = context.InitSession(httpContext.Response.W, httpContext.Request.R)
	receiver.IMiddleware.Invoke(httpContext)
}
