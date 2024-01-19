package minimal

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {

	// 执行过滤器OnActionExecuting
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		traceFiltersDetail := container.Resolve[trace.IManager]().TraceHand("执行过滤器OnActionExecuting：" + parse.ToString(i+1))
		httpContext.Route.Filters[i].OnActionExecuting(httpContext)
		traceFiltersDetail.End(nil)
		// 约定小于1us，不显示
		if traceFiltersDetail.GetTraceDetail().UnTraceTs.Microseconds() <= 1 {
			traceFiltersDetail.Ignore()
		}
	}

	traceDetail := container.Resolve[trace.IManager]().TraceHand("执行路由")
	// 调用action
	callValues := reflect.ValueOf(httpContext.Route.Action).Call(httpContext.Request.Params)
	httpContext.Response.SetValues(callValues...)
	traceDetail.End(nil)

	// 执行过滤器OnActionExecuted
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		traceFiltersDetail := container.Resolve[trace.IManager]().TraceHand("执行过滤器OnActionExecuted" + parse.ToString(i+1))
		httpContext.Route.Filters[i].OnActionExecuted(httpContext)
		traceFiltersDetail.End(nil)
		// 约定小于1us，不显示
		if traceFiltersDetail.GetTraceDetail().UnTraceTs.Microseconds() <= 1 {
			traceFiltersDetail.Ignore()
		}
	}
}
