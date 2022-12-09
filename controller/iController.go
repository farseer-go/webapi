package controller

import (
	"github.com/farseer-go/collections"
	"net/http"
)

// 得到IController接口的所有方法名称
var lstControllerMethodName collections.List[string]

type IController interface {
	init(r *http.Request)         // init 初始化
	getAction() map[string]Action // 获取Action设置
}
