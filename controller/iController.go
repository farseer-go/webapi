package controller

import (
	"github.com/farseer-go/collections"
	"net/http"
	"reflect"
)

// 得到IController接口的所有方法名称
var lstControllerMethodName collections.List[string]

type IController interface {
	// init 初始化
	init(r *http.Request)
}

// 获取IController的方法列表
func getIControllerMethodNames() {
	iControllerType := reflect.TypeOf((*IController)(nil)).Elem()
	lstControllerMethodName = collections.NewList[string]()
	for i := 0; i < iControllerType.NumMethod(); i++ {
		lstControllerMethodName.Add(iControllerType.Method(i).Name)
	}
}
