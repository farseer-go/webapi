package controller

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
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

			// 实例化控制器
			controllerVal := reflect.New(route.controller)
			baseController := getBaseController(controllerVal)

			// 初始化
			baseController.init(r)

			// 入参
			params := baseController.httpContext.GetRequestParam(route.requestParamType, collections.NewList[string]())

			exception.Try(func() {
				// 调用action
				returnVals := controllerVal.MethodByName(route.actionName).Call(params)
				// 初始化返回报文
				baseController.httpContext.InitResponse(returnVals)
				// 输出返回值
				_, _ = w.Write(baseController.httpContext.HttpResponse.BodyBytes)
				// 响应码
				w.WriteHeader(200)
			}).CatchException(func(exp any) {
				// 响应码
				w.WriteHeader(500)
			})
		})
	}
}
