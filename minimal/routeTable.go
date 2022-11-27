package minimal

import (
	"github.com/farseer-go/collections"
	"reflect"
)

var lstRouteTable collections.List[routeTable]

// 路由表
type routeTable struct {
	routeUrl         string                         // 路由地址
	action           any                            // action类型（func）
	requestParamType collections.List[reflect.Type] // 入参
	responseBodyType collections.List[reflect.Type] // 出参
	paramNames       collections.List[string]       // 多字段，方法入参名称
	method           string
}
