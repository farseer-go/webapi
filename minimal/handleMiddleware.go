package minimal

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	// 入参
	params := httpContext.ParseParams()

	sw := stopwatch.StartNew()
	// 执行过滤器OnActionExecuting
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		httpContext.Route.Filters[i].OnActionExecuting(httpContext)
	}

	// 调用action
	httpContext.Response.Body = reflect.ValueOf(httpContext.Route.Action).Call(params)

	// 执行过滤器OnActionExecuted
	for i := 0; i < len(httpContext.Route.Filters); i++ {
		httpContext.Route.Filters[i].OnActionExecuted(httpContext)
	}

	flog.ComponentInfof("webapi", "%s Use：%s", httpContext.URI.Url, sw.GetMillisecondsText())
}
