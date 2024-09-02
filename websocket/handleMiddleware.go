package websocket

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	// 执行过滤器OnActionExecuting
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		httpContext.Route.Filters[i].OnActionExecuting(httpContext)
	}

	// 实现了check.ICheck（必须放在过滤器之后执行）
	httpContext.RequestParamCheck()

	// 调用action
	callValues := reflect.ValueOf(httpContext.Route.Action).Call(httpContext.Request.Params)
	httpContext.Response.SetValues(callValues...)

	// 执行过滤器OnActionExecuted
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		httpContext.Route.Filters[i].OnActionExecuted(httpContext)
	}
}
