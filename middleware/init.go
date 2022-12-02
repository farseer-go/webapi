package middleware

import (
	"github.com/farseer-go/collections"
	"reflect"
)

func Init() {
	MiddlewareList = collections.NewList[IMiddleware](&http{}, &exceptionMiddleware{}, &routing{})
}

// InitMiddleware 初始化管道
func InitMiddleware() {
	for middlewareIndex := 0; middlewareIndex < MiddlewareList.Count(); middlewareIndex++ {
		// 最后一个中间件不需要再设置
		if middlewareIndex+1 == MiddlewareList.Count() {
			return
		}

		curMiddleware := MiddlewareList.Index(middlewareIndex)
		nextMiddleware := MiddlewareList.Index(middlewareIndex + 1)
		SetNextMiddleware(curMiddleware, nextMiddleware)
	}
}

// SetNextMiddleware 设置下一个管道
func SetNextMiddleware(curMiddleware, nextMiddleware IMiddleware) {
	curMiddlewareValue := reflect.ValueOf(curMiddleware)
	// 找到next字段进行赋值下一个中间件管道
	for fieldIndex := 0; fieldIndex < curMiddlewareValue.Elem().NumField(); fieldIndex++ {
		field := curMiddlewareValue.Elem().Field(fieldIndex)
		if field.Type().String() == "middleware.IMiddleware" {
			field.Set(reflect.ValueOf(nextMiddleware))
			break
		}
	}
}
