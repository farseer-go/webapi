package test

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")

	webapi.Area("api/1.0", func() {
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
		webapi.RegisterGET("/mini/hello2", Hello2)
		webapi.RegisterPUT("/mini/hello3", Hello3, "pageSize", "pageIndex")
		webapi.RegisterDELETE("/mini/hello4", Hello4, "pageSize", "pageIndex")
	})
	webapi.RegisterRoutes([]webapi.Route{{Url: "/mini/hello2", Method: "GET", Action: Hello2}})

	webapi.UseCors()
	webapi.UseWebApi()
	webapi.UseStaticFiles()
	webapi.UseApiResponse()

	testPanic(t)

	go webapi.Run("")
	time.Sleep(100 * time.Millisecond)

	testController := TestController{}
	t.Run("api/1.0/test/hello1", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8888/api/1.0/test/hello1", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, testController.Hello1(sizeRequest), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/test/hello2-application/json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8888/api/1.0/test/hello2", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, testController.Hello2(sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/test/hello2-application/json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		val := make(url.Values)
		val.Add("pageSize", strconv.Itoa(sizeRequest.PageSize))
		val.Add("pageIndex", strconv.Itoa(sizeRequest.PageIndex))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/api/1.0/test/hello2", val)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, testController.Hello2(sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/mini/hello1", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8888/api/1.0/mini/hello1", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello1(sizeRequest), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello2", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/mini/hello2")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello2().(pageSizeRequest).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello2().(pageSizeRequest).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
	})
}

func testPanic(t *testing.T) bool {
	return t.Run("TestPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			webapi.RegisterPOST("/", func() {})
		})
	})
}
