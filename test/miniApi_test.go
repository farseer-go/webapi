package test

type pageSizeRequest struct {
	PageSize   int `json:"page_size"`
	PageIndex  int
	noExported string //测试不导出字段
}

func Hello2() any {
	return pageSizeRequest{PageSize: 3, PageIndex: 2}
}
