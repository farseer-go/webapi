package minimal

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
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

	traceDetail := container.Resolve[trace.IManager]().TraceHand("执行路由")
	// 调用action
	callValues := reflect.ValueOf(httpContext.Route.Action).Call(httpContext.Request.Params)
	httpContext.Response.SetValues(callValues...)
	traceDetail.End(nil)

	// 执行过滤器OnActionExecuted
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		httpContext.Route.Filters[i].OnActionExecuted(httpContext)
	}
}
