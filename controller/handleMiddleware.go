package controller

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type HandleMiddleware struct {
}

func (receiver HandleMiddleware) Invoke(httpContext *context.HttpContext) {
	traceDetail := container.Resolve[trace.IManager]().TraceHand("执行路由")
	defer traceDetail.End(nil)

	// 实例化控制器
	controllerVal := reflect.New(httpContext.Route.Controller)

	// 初始化赋控制器
	receiver.initController(httpContext, controllerVal)

	// 自动绑定头部
	receiver.bindHeader(httpContext, controllerVal)

	actionMethod := controllerVal.MethodByName(httpContext.Route.ActionName)

	var callValues []reflect.Value
	// 是否要执行ActionFilter
	if httpContext.Route.IsImplActionFilter {
		actionFilter := controllerVal.Interface().(IActionFilter)
		actionFilter.OnActionExecuting()

		// 实现了check.ICheck（必须放在过滤器之后执行）
		httpContext.RequestParamCheck()
		callValues = actionMethod.Call(httpContext.Request.Params) // 调用action

		actionFilter.OnActionExecuted()
	} else {
		// 实现了check.ICheck（必须放在过滤器之后执行）
		httpContext.RequestParamCheck()
		callValues = actionMethod.Call(httpContext.Request.Params) // 调用action
	}

	httpContext.Response.SetValues(callValues...)
}

// 找到 "controller.BaseController" 字段，并初始化赋值
func (receiver HandleMiddleware) initController(httpContext *context.HttpContext, controllerVal reflect.Value) {
	controllerElem := controllerVal.Elem()
	for i := 0; i < controllerElem.NumField(); i++ {
		fieldVal := controllerElem.Field(i)
		if fieldVal.Type().String() == "controller.BaseController" {
			fieldVal.Set(reflect.ValueOf(BaseController{HttpContext: httpContext}))
			return
		}
	}
}

// 绑定header
func (receiver HandleMiddleware) bindHeader(httpContext *context.HttpContext, controllerVal reflect.Value) {
	// 没有设置绑定字段，则不需要绑定
	if httpContext.Route.AutoBindHeaderName == "" {
		return
	}

	controllerHeaderVal := controllerVal.Elem().FieldByName(httpContext.Route.AutoBindHeaderName)
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
		headerFieldVal.Set(reflect.ValueOf(headerValue))
	}
}
