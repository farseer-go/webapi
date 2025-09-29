package websocket

import (
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"golang.org/x/net/websocket"
)

func SocketHandler(route *context.HttpRoute) websocket.Handler {
	return func(conn *websocket.Conn) {
		// 调用action
		// 解析报文、组装httpContext
		httpContext := context.NewHttpContext(route, nil, conn.Request())
		httpContext.SetWebsocket(conn)

		// InitContext 初始化同一协程上下文，避免在同一协程中多次初始化
		asyncLocal.InitContext()

		// 创建链路追踪上下文
		trackContext := container.Resolve[trace.IManager]().EntryWebSocket(httpContext.URI.Host, httpContext.URI.Url, httpContext.Header.ToMap(), httpContext.URI.GetRealIp())
		trackContext.SetBody(httpContext.Request.BodyString, httpContext.Response.GetHttpCode(), string(httpContext.Response.BodyBytes), nil)
		container.Resolve[trace.IManager]().Push(trackContext, nil)
		//httpContext.Data.Set("Trace", trackContext)

		// 设置到routine，可用于任意子函数获取
		context.RoutineHttpContext.Set(httpContext)
		// 执行第一个中间件
		route.HttpMiddleware.Invoke(httpContext)
		asyncLocal.Release()
	}
}
