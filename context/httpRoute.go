package context

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
)

// HttpRoute 路由表
type HttpRoute struct {
	RouteUrl                string                         // 路由地址
	Controller              reflect.Type                   // 控制器类型
	ControllerName          string                         // 控制器名称
	Action                  any                            // action类型
	ActionName              string                         // action名称
	RequestParamType        collections.List[reflect.Type] // 入参
	ResponseBodyType        collections.List[reflect.Type] // 出参
	Method                  collections.List[string]       // method
	ParamNames              collections.List[string]       // 入参变量名称（显示指定）
	RequestParamIsImplCheck bool                           // 入参DTO是否实现了Check接口
	RequestParamIsModel     bool                           // 入参是否为DTO结构
	ResponseBodyIsModel     bool                           // 出参是否为DTO结构
	AutoBindHeaderName      string                         // 自动绑定header的字段名称
	IsImplActionFilter      bool                           // 是否实现了IActionFilter
	IsGoBasicType           bool                           // 返回值只有一个时，是否为基础类型
	HttpMiddleware          IMiddleware                    // 中间件入口（每个路由的管道都不一样）
	HandleMiddleware        IMiddleware                    // handle中间件
	RouteRegexp             *routeRegexp                   // 正则路由
	Handler                 http.Handler                   // api处理函数
	Filters                 []IFilter                      // 过滤器（对单个路由的执行单元）
	Schema                  string                         // http or ws
}

// JsonToParams json入参转成param
func (receiver *HttpRoute) JsonToParams(request *HttpRequest) []reflect.Value {
	// dto
	if receiver.RequestParamIsModel {
		// 第一个参数，将json反序列化到dto
		firstParamType := receiver.RequestParamType.First() // 先取第一个参数
		val := reflect.New(firstParamType).Interface()
		_ = sonic.Unmarshal(request.BodyBytes, val)
		returnVal := []reflect.Value{reflect.ValueOf(val).Elem()}

		// 第2个参数起，为interface类型，需要做注入操作
		return receiver.parseInterfaceParam(returnVal)
	}

	// 多参数
	return receiver.FormToParams(request.jsonToMap())
}

// FormToParams 将map转成入参值
func (receiver *HttpRoute) FormToParams(mapVal map[string]any) []reflect.Value {
	// dto模式
	if receiver.RequestParamIsModel {
		// 第一个参数，将json反序列化到dto
		firstParamType := receiver.RequestParamType.First() // 先取第一个参数
		// 反序列后的dto对象值
		dtoParamVal := reflect.New(firstParamType).Elem()
		for i := 0; i < firstParamType.NumField(); i++ {
			field := firstParamType.Field(i)
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
				dtoParamVal.FieldByName(field.Name).Set(reflect.ValueOf(fieldVal))
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
			val = receiver.diParam(i)
		} else {
			val = reflect.New(fieldType).Elem().Interface()
			// 指定了参数名称
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
		val := reflect.ValueOf(receiver.diParam(i))
		returnVal = append(returnVal, val)
	}
	return returnVal
}

func (receiver *HttpRoute) diParam(paramIndex int) any {
	// 如果是接口类型，则这里的名称为IocName
	paramName := ""
	if paramIndex < receiver.ParamNames.Count() {
		paramName = receiver.ParamNames.Index(paramIndex)
	}
	iocType := receiver.RequestParamType.Index(paramIndex)
	if !container.IsRegisterType(iocType, paramName) {
		panic(fmt.Sprintf("类型：%s 未注册到IOC", iocType.String()))
	}
	val, _ := container.ResolveType(iocType, paramName)
	return val
}
