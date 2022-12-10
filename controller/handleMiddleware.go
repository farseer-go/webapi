package controller

import (
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	// 实例化控制器
	controllerVal := reflect.New(httpContext.Route.Controller)

	// 初始化赋控制器
	receiver.initController(httpContext, controllerVal)

	// 自动绑定头部
	receiver.bindHeader(httpContext, controllerVal)

	// 入参
	params := httpContext.GetRequestParam()

	// 调用action
	actionMethod := controllerVal.MethodByName(httpContext.Route.ActionName)
	httpContext.Response.Body = actionMethod.Call(params)
}

// 找到 "controller.BaseController" 字段，并初始化赋值
func (receiver HandleMiddleware) initController(httpContext *context.HttpContext, controllerVal reflect.Value) {
	controllerElem := controllerVal.Elem()
	for i := 0; i < controllerElem.NumField(); i++ {
		fieldVal := controllerElem.Field(i)
		if fieldVal.Type().String() == "controller.BaseController" {
			fieldVal.Set(reflect.ValueOf(BaseController{HttpContext: *httpContext}))
			return
		}
	}
}

// 绑定header
func (receiver HandleMiddleware) bindHeader(httpContext *context.HttpContext, controllerVal reflect.Value) {
	controllerElem := controllerVal.Elem()
	controllerType := controllerElem.Type()
	for i := 0; i < controllerElem.NumField(); i++ {
		// 找到需要绑定头部的标记
		if controllerType.Field(i).Tag.Get("webapi") == "header" {
			controllerHeaderVal := controllerElem.Field(i)
			controllerHeaderType := controllerHeaderVal.Type()
			// 遍历需要将header绑定的结构体
			for headerIndex := 0; headerIndex < controllerHeaderVal.NumField(); headerIndex++ {
				headerFieldVal := controllerHeaderVal.Field(headerIndex)
				headerFieldType := headerFieldVal.Type()
				headerName := controllerHeaderType.Field(headerIndex).Tag.Get("webapi")
				if headerName == "" {
					headerName = controllerHeaderType.Field(headerIndex).Name
				}
				headerVal := httpContext.Header.GetValue(headerName)
				if headerVal == "" {
					continue
				}
				headerValue := parse.ConvertValue(headerVal, headerFieldType)
				headerFieldVal.Set(headerValue)
			}
		}
	}
}
