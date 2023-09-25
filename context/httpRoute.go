package context

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"net/http"
	"reflect"
	"strings"
)

// HttpRoute 路由表
type HttpRoute struct {
	RouteUrl            string                         // 路由地址
	Controller          reflect.Type                   // 控制器类型
	ControllerName      string                         // 控制器名称
	Action              any                            // action类型
	ActionName          string                         // action名称
	RequestParamType    collections.List[reflect.Type] // 入参
	ResponseBodyType    collections.List[reflect.Type] // 出参
	Method              collections.List[string]       // method
	ParamNames          collections.List[string]       // 入参变量名称（显示指定）
	RequestParamIsModel bool                           // 是否为DTO结构
	ResponseBodyIsModel bool                           // 是否为DTO结构
	AutoBindHeaderName  string                         // 自动绑定header的字段名称
	IsImplActionFilter  bool                           // 是否实现了IActionFilter
	IsGoBasicType       bool                           // 返回值只有一个时，是否为基础类型
	HttpMiddleware      IMiddleware                    // 中间件入口（每个路由的管道都不一样）
	HandleMiddleware    IMiddleware                    // handle中间件
	RouteRegexp         *routeRegexp                   // 正则路由
	Handler             http.Handler                   // api处理函数
	Filters             []IFilter                      // 过滤器（对单个路由的执行单元）
}

// JsonToParams json入参转成param
func (receiver *HttpRoute) JsonToParams(request *HttpRequest) []reflect.Value {
	// dto
	if receiver.RequestParamIsModel {
		// 第一个参数，将json反序列化到dto
		firstParamType := receiver.RequestParamType.First() // 先取第一个参数
		val := reflect.New(firstParamType).Interface()
		_ = json.Unmarshal(request.BodyBytes, val)
		returnVal := []reflect.Value{reflect.ValueOf(val).Elem()}

		// 第2个参数起，为interface类型，需要做注入操作
		return receiver.parseInterfaceParam(returnVal)
	}

	// 多参数
	mapVal := request.jsonToMap()
	return receiver.FormToParams(mapVal)
}

// FormToParams 将map转成入参值
func (receiver *HttpRoute) FormToParams(mapVal map[string]any) []reflect.Value {
	// dto模式
	if receiver.RequestParamIsModel {
		// 第一个参数，将json反序列化到dto
		dtoParam := receiver.RequestParamType.First()
		// 反序列后的dto对象值
		dtoParamVal := reflect.New(dtoParam).Elem()
		for i := 0; i < dtoParam.NumField(); i++ {
			field := dtoParam.Field(i)
			if !field.IsExported() {
				continue
			}
			// 支持json标签
			key := strings.ToLower(field.Tag.Get("json"))
			if key == "" {
				key = strings.ToLower(field.Name)
			}
			kv, exists := mapVal[key]
			if exists {
				fieldVal := parse.ConvertValue(kv, dtoParamVal.Field(i).Type())
				// dto中的字段赋值
				dtoParamVal.FieldByName(field.Name).Set(fieldVal)
			}
		}
		returnVal := []reflect.Value{dtoParamVal}

		// 第2个参数起，为interface类型，需要做注入操作
		return receiver.parseInterfaceParam(returnVal)
	}

	// 非dto模式
	lstParams := make([]reflect.Value, receiver.RequestParamType.Count())
	for i := 0; i < receiver.RequestParamType.Count(); i++ {
		fieldType := receiver.RequestParamType.Index(i)
		var val any
		// interface类型，则通过注入的方式
		if fieldType.Kind() == reflect.Interface {
			// 如果是接口类型，则这里的名称为IocName
			paramName := ""
			if i < receiver.ParamNames.Count() {
				paramName = receiver.ParamNames.Index(i)
			}
			var err error
			val, err = container.ResolveType(fieldType, paramName)
			if err != nil {
				exception.ThrowWebException(500, err.Error())
			}
		} else {
			val = reflect.New(fieldType).Elem().Interface()
			//
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

// 第2个参数起，为interface类型，需要做注入操作
func (receiver *HttpRoute) parseInterfaceParam(returnVal []reflect.Value) []reflect.Value {
	for i := 1; i < receiver.RequestParamType.Count(); i++ {
		// 如果是接口类型，则这里的名称为IocName
		paramName := ""
		if i < receiver.ParamNames.Count() {
			paramName = receiver.ParamNames.Index(i)
		}
		val, _ := container.ResolveType(receiver.RequestParamType.Index(i), paramName)
		returnVal = append(returnVal, reflect.ValueOf(val))
	}
	return returnVal
}
