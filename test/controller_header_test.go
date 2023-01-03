package test

import (
	"fmt"
	"github.com/farseer-go/webapi/controller"
)

type header struct {
	ContentType  string `webapi:"Content-Type"`
	ContentType2 string
}
type TestHeaderController struct {
	controller.BaseController
	Header header `webapi:"header"`
}

//func (r *TestHeaderController) Base() {
//
//}

func (r *TestHeaderController) Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

func (r *TestHeaderController) Hello2(pageSize int, pageIndex int) pageSizeRequest {
	return pageSizeRequest{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}

func (r *TestHeaderController) Hello3() (TValue string) {
	return r.HttpContext.Header.GetValue("Content-Type")
}

func (r *TestHeaderController) OnActionExecuting() {
	if r.HttpContext.Method != "GET" && r.Header.ContentType == "" {
		panic("测试失败，未获取到：Header.ContentType")
	}
	r.HttpContext.Response.AddHeader("Executing", "true")
	r.HttpContext.Response.SetHeader("Set-Header1", "true")
	r.HttpContext.Response.SetHeader("Set-Header2", "true")
}

func (r *TestHeaderController) OnActionExecuted() {
	r.HttpContext.Response.AddHeader("Executed", "true")
	r.HttpContext.Response.DelHeader("Set-Header2")
}
