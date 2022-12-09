package controller

import (
	"github.com/farseer-go/webapi/context"
	"net/http"
)

type BaseController struct {
	HttpContext context.HttpContext // 上下文
	Action      map[string]Action   // 设置每个Action参数
}

func (receiver BaseController) init(r *http.Request) {

}

func (receiver BaseController) getAction() map[string]Action {
	return receiver.Action
}
