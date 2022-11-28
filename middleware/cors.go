package middleware

import "github.com/farseer-go/webapi/context"

type Cors struct {
	IMiddleware
}

func (receiver *Cors) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)
}
