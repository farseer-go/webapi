package controller

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/webapi/check"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/middleware"
	"reflect"
	"strings"
)

// Register 自动注册控制器下的所有Action方法
func Register(area string, c IController) collections.List[*context.HttpRoute] {
	cVal := reflect.ValueOf(c)
	cType := cVal.Type()
	controllerType := reflect.Indirect(cVal).Type()
	controllerName := strings.TrimSuffix(controllerType.Name(), "Controller")
	actions := c.getAction()

	// 查找是否需要自动绑定header
	autoBindHeaderName := findAutoBindHeaderName(controllerType)

	// 是否实现了IActionFilter
	var actionFilter = reflect.TypeOf((*IActionFilter)(nil)).Elem()
	isImplActionFilter := cType.Implements(actionFilter)

	lst := collections.NewList[*context.HttpRoute]()
	// 遍历controller下的action函数
	for i := 0; i < cType.NumMethod(); i++ {
		actionMethod := cType.Method(i)
		// 如果是来自基类的方法、非导出类型，则跳过
		if actionMethod.IsExported() && !lstControllerMethodName.Contains(actionMethod.Name) {
			httpRoute := registerAction(area, actionMethod, actions, controllerName, controllerType, autoBindHeaderName, isImplActionFilter)
			if httpRoute != nil {
				lst.Add(httpRoute)
			}
		}
	}
	return lst
}

// 查找自动绑定header的字段
func findAutoBindHeaderName(controllerType reflect.Type) string {
	var controllerFieldName string
	for i := 0; i < controllerType.NumField(); i++ {
		// 找到需要绑定头部的标记
		controllerFieldType := controllerType.Field(i)
		if controllerFieldType.Tag.Get("webapi") == "header" {
			controllerFieldName = controllerFieldType.Name
			break
		}
	}
	return controllerFieldName
}

// 注册Action
func registerAction(area string, actionMethod reflect.Method, actions map[string]Action, controllerName string, controllerType reflect.Type, autoBindHeaderName string, isImplActionFilter bool) *context.HttpRoute {
	methodType := actionMethod.Type
	actionName := actionMethod.Name
	// 如果是ActionFilter过滤器，则跳过
	if isImplActionFilter && (actionName == "OnActionExecuting" || actionName == "OnActionExecuted") {
		return nil
	}

	// 控制器都是以方法的形式，因此第0个入参是接收器，应去除
	lstRequestParamType := collections.NewList(types.GetInParam(methodType)...)
	lstRequestParamType.RemoveAt(0)
	lstResponseParamType := collections.NewList(types.GetOutParam(methodType)...)

	// 设置Action默认参数
	if _, exists := actions[actionName]; !exists {
		actions[actionName] = Action{Method: "POST"}
	}

	// 多参数解析
	var paramNames []string
	if actions[actionName].Params != "" {
		paramNames = strings.Split(actions[actionName].Params, ",")
	}

	if actions[actionName].Method == "" {
		action := actions[actionName]
		action.Method = "POST"
		actions[actionName] = action
	}

	// 入参是否为DTO模式
	isDtoModel := types.IsDtoModelIgnoreInterface(lstRequestParamType.ToArray())
	// 是否实现了ICheck
	var requestParamIsImplCheck bool
	if isDtoModel {
		// 是否实现了check.ICheck
		var checker = reflect.TypeOf((*check.ICheck)(nil)).Elem()
		requestParamIsImplCheck = lstRequestParamType.First().Implements(checker)
	}

	// 添加到路由表
	return &context.HttpRoute{
		RouteUrl:                area + strings.ToLower(controllerName) + "/" + strings.ToLower(actionName),
		Action:                  methodType,
		Method:                  collections.NewList(strings.Split(strings.ToUpper(actions[actionName].Method), "|")...),
		RequestParamType:        lstRequestParamType,
		ResponseBodyType:        lstResponseParamType,
		RequestParamIsImplCheck: requestParamIsImplCheck,
		RequestParamIsModel:     isDtoModel,
		ResponseBodyIsModel:     types.IsDtoModel(lstResponseParamType.ToArray()),
		ParamNames:              collections.NewList(paramNames...),
		HttpMiddleware:          &middleware.Http{},
		HandleMiddleware:        &HandleMiddleware{},
		IsGoBasicType:           types.IsGoBasicType(lstResponseParamType.First()),

		AutoBindHeaderName: autoBindHeaderName,
		IsImplActionFilter: isImplActionFilter,
		Controller:         controllerType,
		ControllerName:     controllerName,
		ActionName:         actionName,
	}
}
