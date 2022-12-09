package webapi

import (
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
	"net/http"
)

// 将路由表注册到http.HandleFunc
func handleRoute() {
	// 遍历路由注册表
	for i := 0; i < context.LstRouteTable.Count(); i++ {
		route := context.LstRouteTable.Index(i)
		http.HandleFunc(route.RouteUrl, func(w http.ResponseWriter, r *http.Request) {
			// 组装最后一个API中间件
			lst := middleware.MiddlewareList.ToArray()
			last := lst[middleware.MiddlewareList.Count()-1]

			// minimalApi
			if route.ControllerName == "" && route.ActionName == "" {
				middleware.SetNextMiddleware(last, minimal.HandleMiddleware{})
			} else { // controller
				middleware.SetNextMiddleware(last, controller.HandleMiddleware{})
			}

			// 解析报文、组装httpContext
			httpContext := context.NewHttpContext(route, w, r)

			// 执行第一个中间件
			lst[0].Invoke(&httpContext)
		})
	}
}
