package context

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/vmihailenco/msgpack/v5"
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

	// 依赖注入缓存（在路由全部注册、容器就绪后由 PrecomputeDI 预解析）
	diResolved   bool            // 是否已完成预解析
	diParamCache []reflect.Value // 各入参的预解析单例实例（nil表示该位非缓存项，需按请求解析）
	diParamIsDI  []bool          // 各入参是否为需要注入的interface类型
}

// PrecomputeDI 预解析路由的注入依赖。
// 在所有模块加载、所有路由注册完成后调用一次（非热路径）。
// 对于注册为单例的依赖，提前解析好 reflect.Value 并缓存，
// 请求时直接复用，避免每个请求重复走容器查找、加锁与反射装箱。
func (receiver *HttpRoute) PrecomputeDI() {
	paramCount := receiver.RequestParamType.Count()
	receiver.diParamCache = make([]reflect.Value, paramCount)
	receiver.diParamIsDI = make([]bool, paramCount)

	// 从第2个参数起（或非DTO模式下从首个interface参数起）才是注入参数
	for i := 0; i < paramCount; i++ {
		paramType := receiver.RequestParamType.Index(i)
		if paramType.Kind() != reflect.Interface {
			continue
		}
		receiver.diParamIsDI[i] = true

		// 解析别名（与 diParam 保持一致）
		paramName := ""
		if i < receiver.ParamNames.Count() {
			paramName = receiver.ParamNames.Index(i)
		}
		// 仅单例可安全缓存；临时生命周期保持每请求解析
		if container.IsSingleType(paramType, paramName) {
			if val, err := container.ResolveType(paramType, paramName); err == nil {
				receiver.diParamCache[i] = reflect.ValueOf(val)
			}
		}
	}
	receiver.diResolved = true
}

// JsonToParams json入参转成param
func (receiver *HttpRoute) JsonToParams(request *HttpRequest) []reflect.Value {
	// dto
	if receiver.RequestParamIsModel {
		// 第一个参数，将json反序列化到dto
		firstParamType := receiver.RequestParamType.First() // 先取第一个参数
		val := reflect.New(firstParamType).Interface()
		_ = snc.Unmarshal(request.BodyBytes, val)
		returnVal := []reflect.Value{reflect.ValueOf(val).Elem()}

		// 第2个参数起，为interface类型，需要做注入操作
		return receiver.parseInterfaceParam(returnVal)
	}

	// 多参数
	return receiver.FormToParams(request.jsonToMap())
}

// MsgpackToParams 处理 application/x-msgpack 入参
func (receiver *HttpRoute) MsgpackToParams(request *HttpRequest) []reflect.Value {
	// DTO 模式处理
	if receiver.RequestParamIsModel {
		// 1. 获取第一个参数的类型并创建实例
		firstParamType := receiver.RequestParamType.First()
		val := reflect.New(firstParamType).Interface()

		// 2. 将 msgpack 二进制数据直接反序列化到 DTO 对象中
		err := msgpack.Unmarshal(request.BodyBytes, val)
		if err != nil {
			// 实际业务中建议增加错误处理，例如记录日志
			flog.Errorf("Msgpack Unmarshal 失败: %v", err)
		}

		// 3. 将对象转换为 reflect.Value
		returnVal := []reflect.Value{reflect.ValueOf(val).Elem()}

		// 4. 处理后续注入逻辑
		return receiver.parseInterfaceParam(returnVal)
	}

	// 多参数模式处理：将 msgpack 转为 map 后再处理
	return receiver.FormToParams(request.msgpackToMap())
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
		// interface类型，则通过注入的方式（优先命中预解析单例缓存）
		if fieldType.Kind() == reflect.Interface {
			lstParams[i] = receiver.resolveDIParam(i)
			continue
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
		returnVal = append(returnVal, receiver.resolveDIParam(i))
	}
	return returnVal
}

// resolveDIParam 获取第i个注入参数的值：命中预解析的单例缓存则直接复用，
// 否则（临时生命周期或未预解析）回退到按请求实时解析。
func (receiver *HttpRoute) resolveDIParam(paramIndex int) reflect.Value {
	if receiver.diResolved && paramIndex < len(receiver.diParamCache) {
		if cached := receiver.diParamCache[paramIndex]; cached.IsValid() {
			return cached
		}
	}
	return reflect.ValueOf(receiver.diParam(paramIndex))
}

func (receiver *HttpRoute) diParam(paramIndex int) any {
	// 如果是接口类型，则这里的名称为IocName
	paramName := ""
	if paramIndex < receiver.ParamNames.Count() {
		paramName = receiver.ParamNames.Index(paramIndex)
	}
	iocType := receiver.RequestParamType.Index(paramIndex)
	// 单次解析：未注册时 ResolveType 返回 error，避免额外再做一次 IsRegisterType 的查找与加锁
	val, err := container.ResolveType(iocType, paramName)
	if err != nil {
		panic(fmt.Sprintf("类型：%s 未注册到IOC", iocType.String()))
	}
	return val
}
