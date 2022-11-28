package middleware

import (
	"github.com/farseer-go/webapi/context"
)

type ApiResponse struct {
	IMiddleware
}

func (receiver *ApiResponse) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)
}
