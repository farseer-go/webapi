package controller

import (
	"github.com/farseer-go/collections"
)

// 得到IController接口的所有方法名称
var lstControllerMethodName collections.List[string]

type IController interface {
	getAction() map[string]Action // 获取Action设置
}
