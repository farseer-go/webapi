package webapi

import (
	"github.com/farseer-go/webapi/context"
	"net/http"
)

// 将路由表注册到http.HandleFunc
func handleRoute(mux *http.ServeMux) {
	// 遍历路由注册表
	for i := 0; i < context.LstRouteTable.Count(); i++ {
		route := context.LstRouteTable.Index(i)
		mux.HandleFunc(route.RouteUrl, func(w http.ResponseWriter, r *http.Request) {
			// 解析报文、组装httpContext
			httpContext := context.NewHttpContext(route, w, r)

			// 执行第一个中间件
			route.HttpMiddleware.Invoke(&httpContext)
		})
	}
}
