package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"net/http"
	"testing"
	"time"
)

// BenchmarkRun-12    	    4434	    304151 ns/op	   22731 B/op	     202 allocs/op（优化前）
// BenchmarkRun-12    	    4432	    239972 ns/op	   22758 B/op	     203 allocs/op（第一次优化：HttpResponse.Body改为[]any）
func BenchmarkRun(b *testing.B) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterPOST("/dto", func(req pageSizeRequest) string {
		webapi.GetHttpContext().Response.SetMessage(200, "测试成功")
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})

	webapi.UseApiResponse()
	go webapi.Run(":8094")
	time.Sleep(10 * time.Millisecond)
	b.ReportAllocs()
	sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
	marshal, _ := json.Marshal(sizeRequest)

	for i := 0; i < b.N; i++ {
		rsp, _ := http.Post("http://127.0.0.1:8094/dto", "application/json", bytes.NewReader(marshal))
		_ = rsp.Body.Close()
	}
}
