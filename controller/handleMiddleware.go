package controller

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	// 实例化控制器
	controllerVal := reflect.New(httpContext.Route.Controller)
	baseController := getBaseController(controllerVal)

	// 初始化
	baseController.HttpContext = *httpContext

	// 入参
	params := httpContext.GetRequestParam()

	// 调用action
	httpContext.Response.Body = controllerVal.MethodByName(httpContext.Route.ActionName).Call(params)
}
