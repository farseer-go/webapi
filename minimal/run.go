package minimal

import (
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/context"
	"net/http"
	"reflect"
	"strings"
)

func Run() {
	// 遍历路由注册表
	for i := 0; i < lstRouteTable.Count(); i++ {
		route := lstRouteTable.Index(i)
		http.HandleFunc(route.routeUrl, func(w http.ResponseWriter, r *http.Request) {
			// 检查method
			if strings.ToUpper(route.method) != r.Method {
				// 响应码
				w.WriteHeader(405)
				return
			}

			httpContext := context.NewHttpContext(r)

			// 入参
			params := httpContext.GetRequestParam(route.requestParamType, route.paramNames)

			exception.Try(func() {
				// 调用action
				returnVals := reflect.ValueOf(route.action).Call(params)
				// 初始化返回报文
				httpContext.InitResponse(returnVals)
				// 输出返回值
				_, _ = w.Write(httpContext.HttpResponse.BodyBytes)
				// 响应码
				w.WriteHeader(200)
			}).CatchException(func(exp any) {
				// 响应码
				w.WriteHeader(500)
			})
		})
	}
}
