package controller

import (
	"github.com/farseer-go/collections"
	"reflect"
)

func Init() {
	// 获取IController的方法列表
	getIControllerMethodNames()
}

// 获取IController的方法列表
func getIControllerMethodNames() {
	iControllerType := reflect.TypeOf((*IController)(nil)).Elem()
	lstControllerMethodName = collections.NewList[string]()
	for i := 0; i < iControllerType.NumMethod(); i++ {
		lstControllerMethodName.Add(iControllerType.Method(i).Name)
	}
}
