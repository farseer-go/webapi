package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type psRequest struct {
	PageSize    int `json:"Page_size"`
	PageIndex   int
	noExported  string //测试不导出字段
	CheckResult int
}

// 测试check
func (receiver *psRequest) Check() {
	receiver.CheckResult++
}

func TestRequest(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterPOST("/dto", func(req psRequest) string {
		webapi.GetHttpContext().Response.SetMessage(200, "测试成功")
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d，checkResult=%d", req.PageSize, req.PageIndex, req.CheckResult)
	})

	webapi.RegisterGET("/empty", func() any {
		return psRequest{PageSize: 3, PageIndex: 2}
	})

	webapi.RegisterPUT("/multiParam", func(pageSize int, pageIndex int) psRequest {
		return psRequest{
			PageSize:  pageSize,
			PageIndex: pageIndex,
		}
	}, "page_size", "pageIndex")

	webapi.RegisterPOST("/array", func(ids []int, enable bool) []int {
		return ids
	}, "ids", "enable")
	webapi.UseApiResponse()
	go webapi.Run(":8085")
	time.Sleep(10 * time.Millisecond)

	t.Run("dto-json", func(t *testing.T) {
		sizeRequest := psRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8085/dto", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d，checkResult=1", sizeRequest.PageSize, sizeRequest.PageIndex)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, "测试成功", apiResponse.StatusMessage)
	})

	t.Run("dto-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8085/dto", val)

		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d，checkResult=1", 10, 2)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, "测试成功", apiResponse.StatusMessage)
	})

	t.Run("dto-formData", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8085/dto", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "multipart/form-data")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d，checkResult=1", 10, 2)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, "测试成功", apiResponse.StatusMessage)
	})

	t.Run("empty", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8085/empty")
		apiResponse := core.NewApiResponseByReader[psRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 3, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("multiParam-json", func(t *testing.T) {
		sizeRequest := psRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8085/multiParam", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[psRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, psRequest{PageSize: sizeRequest.PageSize, PageIndex: sizeRequest.PageIndex}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("multiParam-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))

		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8085/multiParam", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[psRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 10, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("multiParam-formData", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))

		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8085/multiParam", strings.NewReader(val.Encode()))
		req.Header.Set("Content-Type", "multipart/form-data")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[psRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 10, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("ids", func(t *testing.T) {
		type array struct {
			Ids    []int
			Enable bool
		}
		b, _ := json.Marshal(array{Ids: []int{1, 2, 3}, Enable: true})
		rsp, _ := http.Post("http://127.0.0.1:8085/array", "application/json", bytes.NewReader(b))
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{1, 2, 3}, apiResponse.Data)
	})
}
