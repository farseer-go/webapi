package test

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
	"testing"
)

func TestApiDoc(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	webapi.UseApiResponse()
	webapi.PrintRoute()
	webapi.UseApiDoc()

	webapi.RegisterPOST("/dto", func(req pageSizeRequest) string {
		webapi.GetHttpContext().Response.SetMessage(200, "测试成功")
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})

	webapi.RegisterGET("/empty", func() any {
		return pageSizeRequest{PageSize: 3, PageIndex: 2}
	})

	webapi.RegisterPUT("/multiParam", func(pageSize int, pageIndex int) pageSizeRequest {
		return pageSizeRequest{
			PageSize:  pageSize,
			PageIndex: pageIndex,
		}
	}, "page_size", "pageIndex")

	webapi.Run("")
}
