package controller

import "github.com/farseer-go/collections"

func Init() {
	// 获取IController的方法列表
	getIControllerMethodNames()
	lstRouteTable = collections.NewList[routeTable]()
}
