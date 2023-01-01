package test

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")

	webapi.Area("/api/1.0/", func() {

		// 自动注册控制器下的所有Action方法
		webapi.RegisterController(&TestController{
			BaseController: controller.BaseController{
				Action: map[string]controller.Action{
					"Hello1": {Method: "POST"},
					"Hello2": {Method: "POST", Params: "pageSize,pageIndex"},
				},
			},
		})

		// 注册单个Api
		webapi.RegisterPOST("/mini/hello1", Hello1)
		webapi.RegisterPOST("/mini/hello2", Hello2)
		webapi.RegisterPOST("/mini/hello3", Hello3, "pageSize", "pageIndex")
		webapi.RegisterPOST("/mini/hello4", Hello4, "pageSize", "pageIndex")
	})

	webapi.UseApiResponse()
	go webapi.Run()
	time.Sleep(1 * time.Second)
}

type pageSizeRequest struct {
	PageSize  int
	PageIndex int
}

func Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

func Hello2(req pageSizeRequest) any {
	return pageSizeRequest{
		PageSize:  req.PageSize,
		PageIndex: req.PageIndex,
	}
}

func Hello3(pageSize int, pageIndex int) pageSizeRequest {
	return pageSizeRequest{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}

func Hello4(pageSize int, pageIndex int) (int, int) {
	return pageSize, pageIndex
}

type TestController struct {
	controller.BaseController
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
