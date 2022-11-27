package controller

import (
	"github.com/farseer-go/collections"
	"reflect"
)

var lstRouteTable collections.List[routeTable]

// 路由表
type routeTable struct {
	routeUrl         string                         // 路由地址
	controller       reflect.Type                   // 控制器类型
	controllerName   string                         // 控制器名称
	action           reflect.Type                   // action类型
	actionName       string                         // action名称
	requestParamType collections.List[reflect.Type] // 入参
	responseBodyType collections.List[reflect.Type] // 出参
	method           string
}
