package context

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
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
	HttpMiddleware      IMiddleware
	HandleMiddleware    IMiddleware
}

// MapToParams 将map转成入参值
func (receiver *HttpRoute) MapToParams(mapVal map[string]any) []reflect.Value {
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
		fieldType := receiver.RequestParamType.Index(i)
		var val any
		// interface类型，则通过注入的方式
		if fieldType.Kind() == reflect.Interface {
			val = container.ResolveType(fieldType)
		} else {
			val = reflect.New(fieldType).Elem().Interface()
			if receiver.ParamNames.Count() > i {
				paramName := strings.ToLower(receiver.ParamNames.Index(i))
				paramVal := mapVal[paramName]
				val = parse.Convert(paramVal, val)

				// 当实际只有一个接收参数时，不需要指定参数
			} else if receiver.ParamNames.Count() == 0 && len(mapVal) == 1 {
				for _, paramVal := range mapVal {
					val = parse.Convert(paramVal, val)
				}
			}
		}

		lstParams[i] = reflect.ValueOf(val)
	}
	return lstParams
}
