package controller

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type ControllerMiddleware struct {
}

func (receiver ControllerMiddleware) Invoke(httpContext *context.HttpContext) {
	// 实例化控制器
	controllerVal := reflect.New(httpContext.HttpRoute.Controller)
	baseController := getBaseController(controllerVal)

	// 初始化
	baseController.httpContext = *httpContext

	// 入参
	params := httpContext.GetRequestParam()

	// 调用action
	httpContext.HttpResponse.Body = controllerVal.MethodByName(httpContext.HttpRoute.ActionName).Call(params)
}
