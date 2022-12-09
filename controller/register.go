package controller

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/webapi/context"
	"reflect"
	"strings"
)

// Register 自动注册控制器下的所有Action方法
func Register(area string, c IController) {
	cVal := reflect.ValueOf(c)
	cType := cVal.Type()
	cRealType := reflect.Indirect(cVal).Type()
	controllerName := strings.TrimSuffix(cRealType.Name(), "Controller")
	actions := c.getAction()

	// 遍历controller下的action函数
	for i := 0; i < cType.NumMethod(); i++ {
		methodType := cType.Method(i).Type
		actionName := cType.Method(i).Name
		// 如果是来自基类的方法、非导出类型，则跳过
		if !cType.Method(i).IsExported() || lstControllerMethodName.Contains(actionName) {
			continue
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
		}

		// 添加到路由表
		context.LstRouteTable.Add(context.HttpRoute{
			RouteUrl:            area + strings.ToLower(controllerName) + "/" + strings.ToLower(actionName),
			Controller:          cRealType,
			ControllerName:      controllerName,
			Action:              methodType,
			ActionName:          actionName,
			RequestParamType:    lstRequestParamType,
			ResponseBodyType:    lstResponseParamType,
			RequestParamIsModel: types.IsDtoModel(lstRequestParamType.ToArray()),
			ResponseBodyIsModel: types.IsDtoModel(lstResponseParamType.ToArray()),
			Method:              strings.ToUpper(actions[actionName].Method),
			ParamNames:          collections.NewList(paramNames...),
		})
	}
}

// InSlice checks given string in string slice or not.
func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
