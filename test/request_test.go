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

func TestRequest(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterPOST("/dto", func(req pageSizeRequest) string {
		return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
	})

	webapi.RegisterGET("/empty", func() any {
		return pageSizeRequest{PageSize: 3, PageIndex: 2}
	})

	webapi.RegisterPUT("/multiParam", func(pageSize int, pageIndex int) pageSizeRequest {
		return pageSizeRequest{
			PageSize:  pageSize,
			PageIndex: pageIndex,
		}
	}, "page_size", "pageIndex")

	webapi.UseApiResponse()
	go webapi.Run(":8085")
	time.Sleep(10 * time.Millisecond)

	t.Run("dto-json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8085/dto", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", sizeRequest.PageSize, sizeRequest.PageIndex)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("dto-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8085/dto", val)

		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", 10, 2)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
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

		expected := fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", 10, 2)
		assert.Equal(t, expected, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("empty", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8085/empty")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 3, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("multiParam-json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8085/multiParam", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, pageSizeRequest{PageSize: sizeRequest.PageSize, PageIndex: sizeRequest.PageIndex}, apiResponse.Data)
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
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
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
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 10, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})
}
