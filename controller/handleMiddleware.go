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
	controllerElem := controllerVal.Elem()
	// 找到 "controller.BaseController" 字段，并初始化赋值
	for i := 0; i < controllerElem.NumField(); i++ {
		fieldVal := controllerElem.Field(i)
		if fieldVal.Type().String() == "controller.BaseController" {
			fieldVal.Set(reflect.ValueOf(BaseController{HttpContext: *httpContext}))
		}
	}

	// 入参
	params := httpContext.GetRequestParam()

	// 调用action
	actionMethod := controllerVal.MethodByName(httpContext.Route.ActionName)
	httpContext.Response.Body = actionMethod.Call(params)
}
