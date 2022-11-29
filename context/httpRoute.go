package context

import (
	"github.com/farseer-go/collections"
	"reflect"
)

var LstRouteTable collections.List[HttpRoute]

// HttpRoute 路由表
type HttpRoute struct {
	RouteUrl            string                         // 路由地址
	Controller          reflect.Type                   // 控制器类型
	ControllerName      string                         // 控制器名称
	Action              any                            // action类型
	ActionName          string                         // action名称
	RequestParamType    collections.List[reflect.Type] // 入参
	ResponseBodyType    collections.List[reflect.Type] // 出参
	Method              string
	ParamNames          collections.List[string]
	RequestParamIsModel bool // 是否为DTO结构
	ResponseBodyIsModel bool // 是否为DTO结构
}
