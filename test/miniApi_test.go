package test

import (
	"fmt"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/action"
)

type pageSizeRequest struct {
	PageSize   int `json:"page_size"`
	PageIndex  int
	noExported string //测试不导出字段
}

func Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

func Hello2() any {
	return pageSizeRequest{PageSize: 3, PageIndex: 2}
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

func Hello5() {
	exception.ThrowWebException(501, "s501")
}

func Hello6() {
	exception.ThrowException("s500")
}

func Hello7(actionType int, testInject ITestInject) action.IResult {
	if testInject.Call() != "ok" {
		panic("inject error")
	}

	switch actionType {
	case 0:
		return action.Redirect("/api/1.0/mini/hello2")
	case 1:
		return action.View("")
	case 2:
		return action.View("mini/hello7")
	case 3:
		return action.View("mini/hello7.txt")
	case 4:
		return action.Content("ccc")
	case 5:
		return action.FileContent("./views/mini/hello7.log")
	}

	return action.Content("eee")
}

func Hello8() {
}

func Hello9(req pageSizeRequest, testInject ITestInject) string {
	if testInject.Call() != "ok" {
		panic("inject error")
	}
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}
