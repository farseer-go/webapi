package test

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.Area("api/1.0", func() {
		// 自动注册控制器下的所有Action方法
		webapi.RegisterController(&TestController{
			BaseController: controller.BaseController{
				Action: map[string]controller.Action{
					"Hello1": {Method: "POST"},
					"Hello2": {Method: "POST", Params: "pageSize,pageIndex"},
					"Hello3": {Method: "GET"},
				},
			},
		})

		// 注册单个Api
		webapi.RegisterPOST("/mini/hello1", Hello1)
		webapi.RegisterGET("/mini/hello2", Hello2)
		webapi.RegisterPUT("/mini/hello3", Hello3, "pageSize", "pageIndex")
		webapi.RegisterDELETE("/mini/hello4", Hello4, "pageSize", "pageIndex")
	})
	webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello2", Method: "GET", Action: Hello2})
	webapi.RegisterPOST("/mini/hello5", Hello5)
	webapi.RegisterPOST("/mini/hello6", Hello6)
	webapi.RegisterPOST("/mini/hello7", Hello7)
	assert.Panics(t, func() {
		webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello3", Method: "GET", Action: Hello2, Params: []string{"aaa"}})
	})
	webapi.UseCors()
	webapi.UseWebApi()
	webapi.UseStaticFiles()
	webapi.UseApiResponse()
	webapi.RegisterMiddleware(&middleware.UrlRewriting{})

	testPanic(t)

	go webapi.Run("")
	time.Sleep(10 * time.Millisecond)

	server := webapi.NewApplicationBuilder()
	server.RegisterPOST("/mini/hello5", Hello5)
	server.RegisterPOST("/mini/hello6", Hello6)
	go server.Run(":8889")
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

	t.Run("api/1.0/test/hello1-GET", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/api/1.0/test/hello1")
		assert.Equal(t, 405, rsp.StatusCode)
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

	t.Run("api/1.0/test/hello2-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("pageSize", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/api/1.0/test/hello2", val)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, testController.Hello2(10, 2), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/test/hello3", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/api/1.0/test/hello3")
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "", apiResponse.Data)
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

	t.Run("api/1.0/mini/hello2", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/api/1.0/mini/hello2")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello2().(pageSizeRequest).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello2().(pageSizeRequest).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/mini/hello3-application/json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8888/api/1.0/mini/hello3", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello3(sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("api/1.0/mini/hello3-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("pageSize", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))

		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8888/api/1.0/mini/hello3", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello3(10, 2), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

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
		assert.Equal(t, Hello2().(pageSizeRequest).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello2().(pageSizeRequest).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello5", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8888/mini/hello5", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[any](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 501, apiResponse.StatusCode)
		assert.Equal(t, "s501", apiResponse.StatusMessage)
		assert.Equal(t, false, apiResponse.Status)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello6", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8888/mini/hello6", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[any](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 500, apiResponse.StatusCode)
		assert.Equal(t, "s500", apiResponse.StatusMessage)
		assert.Equal(t, false, apiResponse.Status)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-0", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(0))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello2().(pageSizeRequest).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello2().(pageSizeRequest).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-1", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(1))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "aaa", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-2", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "aaa", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-3", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(3))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "bbb", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-4", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(4))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ccc", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7-5", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(5))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ddd", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello7--1", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(-1))
		rsp, _ := http.PostForm("http://127.0.0.1:8888/mini/hello7", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "eee", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello5:8889", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello5", "application/json", nil)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "s501", string(body))
		assert.Equal(t, 501, rsp.StatusCode)
	})

	t.Run("mini/hello6:8889", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello6", "application/json", nil)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "s500", string(body))
		assert.Equal(t, 500, rsp.StatusCode)
	})
}

func testPanic(t *testing.T) bool {
	return t.Run("TestPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			webapi.RegisterPOST("/", func() {})
		})
	})
}
