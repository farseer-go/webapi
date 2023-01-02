package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

// MiddlewareList 已注册的中间件集合
var MiddlewareList collections.List[context.IMiddleware]

func Init() {
	MiddlewareList = collections.NewList[context.IMiddleware]()
}

// InitMiddleware 初始化管道
func InitMiddleware() {
	for i := 0; i < context.LstRouteTable.Count(); i++ {
		route := context.LstRouteTable.Index(i)
		// 组装管道
		lstMiddleware := assembledPipeline(route)
		for middlewareIndex := 0; middlewareIndex < lstMiddleware.Count(); middlewareIndex++ {
			// 最后一个中间件不需要再设置
			if middlewareIndex+1 == lstMiddleware.Count() {
				break
			}
			curMiddleware := lstMiddleware.Index(middlewareIndex)
			nextMiddleware := lstMiddleware.Index(middlewareIndex + 1)
			SetNextMiddleware(curMiddleware, nextMiddleware)
		}
	}
}

// 组装管道
func assembledPipeline(route context.HttpRoute) collections.List[context.IMiddleware] {
	// 添加系统中间件
	lst := collections.NewList[context.IMiddleware](route.HttpMiddleware, &exceptionMiddleware{}, &routing{})
	// 添加用户自定义中间件
	for i := 0; i < MiddlewareList.Count(); i++ {
		valIns := reflect.New(reflect.TypeOf(MiddlewareList.Index(i)).Elem()).Interface()
		lst.Add(valIns.(context.IMiddleware))
	}
	// 添加Handle中间件
	lst.Add(route.HandleMiddleware)
	return lst
}

// SetNextMiddleware 设置下一个管道
func SetNextMiddleware(curMiddleware, nextMiddleware context.IMiddleware) {
	curMiddlewareValue := reflect.ValueOf(curMiddleware)
	// 找到next字段进行赋值下一个中间件管道
	for fieldIndex := 0; fieldIndex < curMiddlewareValue.Elem().NumField(); fieldIndex++ {
		field := curMiddlewareValue.Elem().Field(fieldIndex)
		if field.Type().String() == "context.IMiddleware" {
			field.Set(reflect.ValueOf(nextMiddleware))
			break
		}
	}
}

// AddMiddleware 添加中间件
func AddMiddleware(m context.IMiddleware) {
	MiddlewareList.Add(m)
}
