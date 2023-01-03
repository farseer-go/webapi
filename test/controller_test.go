package test

import (
	"fmt"
	"github.com/farseer-go/webapi/controller"
)

type TestController struct {
	controller.BaseController
}

func (r *TestController) Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

// 不导出
func (r *TestController) hello2(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}
