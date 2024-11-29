package test

import (
	"bytes"

	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterDELETE("/multiResponse", func(pageSize int, pageIndex int) (int, int) {
		return pageSize, pageIndex
	}, "page_size", "pageIndex")

	webapi.RegisterPOST("/basicTypeResponse", func(req pageSizeRequest) string {
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})

	webapi.UseApiResponse()
	go webapi.Run(":8086")
	time.Sleep(10 * time.Millisecond)

	t.Run("multiResponse", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := snc.Marshal(sizeRequest)
		req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8086/multiResponse", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{10, 2}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("basicTypeResponse", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := snc.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8086/basicTypeResponse", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})
}
