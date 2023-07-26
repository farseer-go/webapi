package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/action"
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
		// 注册单个Api
		webapi.RegisterPOST("/mini/hello1", Hello1)
		webapi.RegisterGET("/mini/hello2", Hello2)
		webapi.RegisterPUT("/mini/hello3", Hello3, "page_size", "pageIndex")
		webapi.RegisterDELETE("/mini/hello4", Hello4, "page_size", "pageIndex")
	})
	webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello2", Method: "GET", Action: Hello2})
	webapi.RegisterPOST("/mini/hello7", func(actionType int) action.IResult {
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
	})
	webapi.RegisterPOST("/mini/hello9", func(req pageSizeRequest) string {
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})
	webapi.RegisterPOST("/mini/hello10", func() string {
		return webapi.GetHttpContext().ContentType
	})
	webapi.RegisterGET("/mini/hello4/{pageSize}-{pageIndex}", Hello4)
	webapi.RegisterPOST("/mini/hello4/{pageSize}/{pageIndex}", Hello4)

	assert.Panics(t, func() {
		webapi.RegisterRoutes(webapi.Route{Url: "/mini/hello3", Method: "GET", Action: Hello2, Params: []string{"aaa"}})
	})
	webapi.UseCors()
	webapi.UseWebApi()
	webapi.UseStaticFiles()
	webapi.UseApiResponse()
	webapi.RegisterMiddleware(&middleware.UrlRewriting{})

	go webapi.Run("")
	time.Sleep(10 * time.Millisecond)

	t.Run("api/1.0/mini/hello1", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8888/api/1.0/mini/hello1", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello1(sizeRequest), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("api/1.0/mini/hello2", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8888/api/1.0/mini/hello2")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello2().(pageSizeRequest).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello2().(pageSizeRequest).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
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
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("api/1.0/mini/hello3-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))

		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8888/api/1.0/mini/hello3", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, Hello3(10, 2).PageSize, apiResponse.Data.PageSize)
		assert.Equal(t, Hello3(10, 2).PageIndex, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
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
		assert.Equal(t, 200, apiResponse.StatusCode)
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
		assert.Equal(t, 200, apiResponse.StatusCode)
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

	t.Run("/mini/hello4/{pageSize}-{pageIndex}-get", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8888/mini/hello4/15-6", nil)
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{15, 6}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("/mini/hello4/{pageSize}/{pageIndex}-post", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/mini/hello4/15/6", nil)
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{15, 6}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
