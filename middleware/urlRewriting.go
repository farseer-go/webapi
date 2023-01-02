package middleware

import "github.com/farseer-go/webapi/context"

type UrlRewriting struct {
	context.IMiddleware
}

func (receiver *UrlRewriting) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)
}
