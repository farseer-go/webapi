package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/action"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestActionResult(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.RegisterGET("/redirect", func() any {
		return pageSizeRequest{PageSize: 3, PageIndex: 2}
	})

	webapi.RegisterPOST("/mini/testActionResult", func(actionType int) action.IResult {
		switch actionType {
		case 0:
			return action.Redirect("/redirect")
		case 1:
			return action.View("")
		case 2:
			return action.View("testActionResult")
		case 3:
			return action.View("test.txt")
		case 4:
			return action.Content("ccc")
		case 5:
			return action.FileContent("./views/mini/test.log")
		}

		return action.Content("eee")
	})
	webapi.UseApiResponse()

	go webapi.Run(":8088")
	time.Sleep(10 * time.Millisecond)

	t.Run("test-0", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(0))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 3, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("test-1", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(1))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "aaa", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("test-2", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "aaa", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("test-3", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(3))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "bbb", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("test-4", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(4))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ccc", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("test-5", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(5))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ddd", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("test--1", func(t *testing.T) {
		val := make(url.Values)
		val.Add("actionType", strconv.Itoa(-1))
		rsp, _ := http.PostForm("http://127.0.0.1:8088/mini/testActionResult", val)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "eee", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
