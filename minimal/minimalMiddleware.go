package minimal

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type MinimalMiddleware struct {
}

func (receiver MinimalMiddleware) Invoke(httpContext *context.HttpContext) {
	// 入参
	params := httpContext.GetRequestParam()
	// 调用action
	httpContext.HttpResponse.Body = reflect.ValueOf(httpContext.HttpRoute.Action).Call(params)
}
