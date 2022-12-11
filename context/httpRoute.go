package context

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"reflect"
	"strings"
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
	RequestParamIsModel bool   // 是否为DTO结构
	ResponseBodyIsModel bool   // 是否为DTO结构
	AutoBindHeaderName  string // 自动绑定header的字段名称
	IsImplActionFilter  bool   // 是否实现了IActionFilter
}

// 将map转成入参值
func (receiver *HttpRoute) mapToParams(mapVal map[string]any) []reflect.Value {
	// dto模式
	if receiver.RequestParamIsModel {
		param := receiver.RequestParamType.First()
		paramVal := reflect.New(param).Elem()
		for i := 0; i < param.NumField(); i++ {
			field := param.Field(i)
			if !field.IsExported() {
				continue
			}
			key := strings.ToLower(field.Name)
			kv, exists := mapVal[key]
			if exists {
				defVal := paramVal.Field(i).Interface()
				paramVal.FieldByName(field.Name).Set(reflect.ValueOf(parse.Convert(kv, defVal)))
			}
		}
		return []reflect.Value{paramVal}
	}

	// 多参数
	lstParams := make([]reflect.Value, receiver.RequestParamType.Count())
	for i := 0; i < receiver.RequestParamType.Count(); i++ {
		defVal := reflect.New(receiver.RequestParamType.Index(i)).Elem().Interface()
		if receiver.ParamNames.Count() > i {
			paramName := strings.ToLower(receiver.ParamNames.Index(i))
			paramVal, _ := mapVal[paramName]
			defVal = parse.Convert(paramVal, defVal)
		}
		lstParams[i] = reflect.ValueOf(defVal)
	}
	return lstParams
}
