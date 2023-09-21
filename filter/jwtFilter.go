package filter

import (
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
)

type JwtFilter struct {
}

func (receiver *JwtFilter) OnActionExecuting(httpContext *context.HttpContext) {
	if !httpContext.Jwt.Valid() {
		exception.ThrowWebExceptionf(context.InvalidStatusCode, context.InvalidMessage)
	}
}

func (receiver *JwtFilter) OnActionExecuted(httpContext *context.HttpContext) {
}
