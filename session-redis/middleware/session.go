package middleware

import (
	webapiContext "github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/session-redis/context"
)

type Session struct {
	webapiContext.IMiddleware
}

func (receiver *Session) Invoke(httpContext *webapiContext.HttpContext) {
	httpContext.Session = context.InitSession(httpContext.Response.W, httpContext.Request.R)
	receiver.IMiddleware.Invoke(httpContext)
}
