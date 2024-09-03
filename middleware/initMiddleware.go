package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
	"reflect"
)

// InitMiddleware 初始化管道
func InitMiddleware(lstRouteTable map[string]*context.HttpRoute, lstMiddleware collections.List[context.IMiddleware]) {
	for _, route := range lstRouteTable {
		// 系统web是没有中间件的
		if route.HttpMiddleware == nil {
			continue
		}
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
func assembledPipeline(route *context.HttpRoute, lstMiddleware collections.List[context.IMiddleware]) collections.List[context.IMiddleware] {
	// 添加系统中间件
	lst := collections.NewList[context.IMiddleware](route.HttpMiddleware, &exceptionMiddleware{}, &routing{})

	// 添加用户自定义中间件
	for i := 0; i < lstMiddleware.Count(); i++ {
		middlewareType := reflect.TypeOf(lstMiddleware.Index(i)).Elem()
		if route.Schema != "ws" || middlewareType.String() != "middleware.ApiResponse" {
			valIns := reflect.New(middlewareType).Interface()
			lst.Add(valIns.(context.IMiddleware))
		}
	}

	// 添加Handle中间件
	lst.Add(route.HandleMiddleware)

	// 找到cors中间件，放入到http之后（即移到第2个索引）
	for i := 3; i < lst.Count(); i++ {
		if corsMiddleware, isHaveCors := lst.Index(i).(*Cors); isHaveCors {
			lst.RemoveAt(i)
			lst.Insert(1, corsMiddleware)
			break
		}
	}
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
