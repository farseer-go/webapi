package controller

import (
	"github.com/farseer-go/collections"
	"net/http"
	"reflect"
)

func Run() {
	// 遍历路由注册表
	for i := 0; i < lstRouteTable.Count(); i++ {
		route := lstRouteTable.Index(i)
		http.HandleFunc(route.routeUrl, func(w http.ResponseWriter, r *http.Request) {
			// 实例化控制器
			controllerVal := reflect.New(route.controller)
			baseController := getBaseController(controllerVal)

			// 初始化
			baseController.init(r)

			// 入参
			params := baseController.httpContext.GetRequestParam(route.requestParamType, collections.NewList[string]())

			// 调用action
			returnVals := controllerVal.MethodByName(route.actionName).Call(params)

			// 初始化返回报文
			baseController.httpContext.InitResponse(returnVals)

			// 输出返回值
			w.Write(baseController.httpContext.HttpResponse.BodyBytes)

			// 响应码
			w.WriteHeader(200)
		})
	}
}
