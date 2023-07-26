package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.Area("api/1.0", func() {
		webapi.RegisterDELETE("/mini/hello4", func(pageSize int, pageIndex int) (int, int) {
			return pageSize, pageIndex
		}, "page_size", "pageIndex")
	})
	webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello2", Method: "GET", Action: Hello2})

	webapi.RegisterPOST("/mini/hello9", func(req pageSizeRequest) string {
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})
	webapi.RegisterPOST("/mini/hello10", func() string {
		return webapi.GetHttpContext().ContentType
	})

	assert.Panics(t, func() {
		webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello3", Method: "GET", Action: Hello2, Params: []string{"aaa"}})
	})
	webapi.UseWebApi()
	webapi.UseStaticFiles()
	webapi.UseApiResponse()
	webapi.RegisterMiddleware(&middleware.UrlRewriting{})

	go webapi.Run("")
	time.Sleep(10 * time.Millisecond)

	t.Run("api/1.0/mini/hello4", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8888/api/1.0/mini/hello4", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{10, 2}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello2", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/mini/hello2")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 3, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("mini/hello9", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8888/mini/hello9", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello1(sizeRequest), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("mini/hello10", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8888/mini/hello10", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "application/json", apiResponse.Data)
	})
}
