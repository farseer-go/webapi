package controller

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi/context"
	"net/http"
)

type BaseController struct {
	httpContext context.HttpContext
}

func (c *BaseController) init(r *http.Request) {
	flog.Debug("完成初始化")
}
