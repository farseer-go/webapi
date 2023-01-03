package test

import (
	"fmt"
	"github.com/farseer-go/webapi/controller"
)

type TestController struct {
	controller.BaseController
	Header struct{} `webapi:"header"`
}

func (r *TestController) Base() {

}

func (r *TestController) Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

func (r *TestController) Hello2(pageSize int, pageIndex int) pageSizeRequest {
	return pageSizeRequest{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}

func (r *TestController) Hello3() (TValue string) {
	return r.HttpContext.Header.GetValue("Content-Type")
}

func (r *TestController) OnActionExecuting() {
	r.HttpContext.Response.AddHeader("Executing", "true")
}

func (r *TestController) OnActionExecuted() {
	r.HttpContext.Response.AddHeader("Executed", "true")
}
