package test

import (
	"fmt"
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
