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
	// 调用action
	httpContext.Response.Body = reflect.ValueOf(httpContext.Route.Action).Call(params)
	flog.ComponentInfof("webapi", "%s Use：%s", httpContext.URI.Url, sw.GetMillisecondsText())
}
