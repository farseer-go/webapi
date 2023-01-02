package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

// InitMiddleware 初始化管道
func InitMiddleware(lstRouteTable collections.List[context.HttpRoute], lstMiddleware collections.List[context.IMiddleware]) {
	for i := 0; i < lstRouteTable.Count(); i++ {
		route := lstRouteTable.Index(i)
		// 组装管道
		lstPipeline := assembledPipeline(route, lstMiddleware)
		for middlewareIndex := 0; middlewareIndex < lstPipeline.Count(); middlewareIndex++ {
			// 最后一个中间件不需要再设置
			if middlewareIndex+1 == lstPipeline.Count() {
				break
			}
			curMiddleware := lstPipeline.Index(middlewareIndex)
			nextMiddleware := lstPipeline.Index(middlewareIndex + 1)
			setNextMiddleware(curMiddleware, nextMiddleware)
		}
	}
}

// 组装管道
func assembledPipeline(route context.HttpRoute, lstMiddleware collections.List[context.IMiddleware]) collections.List[context.IMiddleware] {
	// 添加系统中间件
	lst := collections.NewList[context.IMiddleware](route.HttpMiddleware, &exceptionMiddleware{}, &routing{})
	// 添加用户自定义中间件
	for i := 0; i < lstMiddleware.Count(); i++ {
		valIns := reflect.New(reflect.TypeOf(lstMiddleware.Index(i)).Elem()).Interface()
		lst.Add(valIns.(context.IMiddleware))
	}
	// 添加Handle中间件
	lst.Add(route.HandleMiddleware)
	return lst
}

// setNextMiddleware 设置下一个管道
func setNextMiddleware(curMiddleware, nextMiddleware context.IMiddleware) {
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
