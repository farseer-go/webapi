package minimal

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	// 入参
	params := httpContext.BuildActionInValue()
	// 调用action
	httpContext.Response.Body = reflect.ValueOf(httpContext.Route.Action).Call(params)
}
