package controller

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// Register 自动注册控制器下的所有Action方法
func Register(area string, c IController) {
	cVal := reflect.ValueOf(c)
	cType := cVal.Type()
	cRealType := reflect.Indirect(cVal).Type()
	controllerName := strings.TrimSuffix(cRealType.Name(), "Controller")

	for i := 0; i < cType.NumMethod(); i++ {
		methodType := cType.Method(i).Type
		actionName := cType.Method(i).Name
		// 如果是来自基类的方法、非导出类型，则跳过
		if !cType.Method(i).IsExported() || lstControllerMethodName.Contains(actionName) {
			continue
		}

		// 控制器都是以方法的形式，因此第0个入参是接收器，应去除
		lstParamType := collections.NewList(types.GetInParam(methodType)...)
		lstParamType.RemoveAt(0)

		// 添加到路由表
		lstRouteTable.Add(routeTable{
			routeUrl:         area + controllerName + "/" + actionName,
			controller:       cRealType,
			controllerName:   controllerName,
			action:           methodType,
			actionName:       actionName,
			requestParamType: lstParamType,
			responseBodyType: collections.NewList(types.GetOutParam(methodType)...),
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
