package webapi

import (
	"net/http"

	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
)

func HttpHandler(route *context.HttpRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// InitContext 初始化同一协程上下文，避免在同一协程中多次初始化
		asyncLocal.InitContext()
		// 解析报文、组装httpContext
		httpContext := context.NewHttpContext(route, w, r)
		// 创建链路追踪上下文
		trackContext := container.Resolve[trace.IManager]().EntryWebApi(httpContext.URI.Host, httpContext.URI.Url, httpContext.Method, httpContext.ContentType, httpContext.Header.ToMap(), httpContext.URI.GetRealIp())
		// 记录出入参
		defer func() {
			trackContext.SetBody(httpContext.Request.BodyString, httpContext.Response.GetHttpCode(), string(httpContext.Response.BodyBytes), httpContext.ResponseHeader.ToMap())
			container.Resolve[trace.IManager]().Push(trackContext, nil)
		}()
		httpContext.Data.Set("Trace", trackContext)

		// 设置到routine，可用于任意子函数获取
		context.RoutineHttpContext.Set(httpContext)
		// 执行第一个中间件
		route.HttpMiddleware.Invoke(httpContext)
		asyncLocal.Release()
	}
}
