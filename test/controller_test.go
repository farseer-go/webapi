package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type header struct {
	ContentType  string `webapi:"Content-Type"`
	ContentType2 string
}
type TestHeaderController struct {
	controller.BaseController
	Header header `webapi:"header"`
}

func (r *TestHeaderController) Hello1(req pageSizeRequest) string {
	r.Response.SetMessage("修改成功")
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

func (r *TestHeaderController) Hello2(pageSize int, pageIndex int) pageSizeRequest {
	return pageSizeRequest{
		PageSize:  pageSize,
		PageIndex: pageIndex,
	}
}

func (r *TestHeaderController) Hello3() (TValue string) {
	return r.HttpContext.Header.GetValue("Content-Type")
}

func (r *TestHeaderController) OnActionExecuting() {
	if r.HttpContext.Method != "GET" && r.Header.ContentType == "" {
		panic("测试失败，未获取到：Header.ContentType")
	}
	r.HttpContext.Response.AddHeader("Executing", "true")
	r.HttpContext.Response.SetHeader("Set-Header1", "true")
	r.HttpContext.Response.SetHeader("Set-Header2", "true")
}

func (r *TestHeaderController) OnActionExecuted() {
	r.HttpContext.Response.AddHeader("Executed", "true")
	r.HttpContext.Response.DelHeader("Set-Header2")
}

func TestController(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.Area("api/1.0", func() {
		// 自动注册控制器下的所有Action方法
		webapi.RegisterController(&TestHeaderController{
			BaseController: controller.BaseController{
				Action: map[string]controller.Action{
					"Hello1": {Method: ""},
					"Hello2": {Method: "POST", Params: "page_Size,pageIndex"},
					"Hello3": {Method: "GET"},
				},
			},
		})
	})

	webapi.UseCors()
	webapi.UseWebApi()
	webapi.UseStaticFiles()
	webapi.UseApiResponse()
	webapi.RegisterMiddleware(&middleware.UrlRewriting{})

	go webapi.Run(":8079")
	time.Sleep(10 * time.Millisecond)

	t.Run("api/1.0/test/hello1", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8079/api/1.0/testheader/hello1", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, "true", rsp.Header.Get("Executing"))
		assert.Equal(t, "true", rsp.Header.Get("Executed"))
		assert.Equal(t, "true", rsp.Header.Get("Set-Header1"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header2"))
	})

	t.Run("api/1.0/test/hello1-GET", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8079/api/1.0/testheader/hello1")
		assert.Equal(t, 405, rsp.StatusCode)
		assert.Equal(t, "", rsp.Header.Get("Executing"))
		assert.Equal(t, "", rsp.Header.Get("Executed"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header1"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header2"))
	})

	t.Run("api/1.0/test/hello2-application/json", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8079/api/1.0/testheader/hello2", "application/json", bytes.NewReader(marshal))
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		controller := TestHeaderController{}
		assert.Equal(t, controller.Hello2(sizeRequest.PageSize, sizeRequest.PageIndex), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, "true", rsp.Header.Get("Executing"))
		assert.Equal(t, "true", rsp.Header.Get("Executed"))
		assert.Equal(t, "true", rsp.Header.Get("Set-Header1"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header2"))
	})

	t.Run("api/1.0/test/hello2-form", func(t *testing.T) {
		val := make(url.Values)
		val.Add("page_Size", strconv.Itoa(10))
		val.Add("pageIndex", strconv.Itoa(2))
		rsp, _ := http.PostForm("http://127.0.0.1:8079/api/1.0/testheader/hello2", val)
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		controller := TestHeaderController{}
		assert.Equal(t, controller.Hello2(10, 2), apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, "true", rsp.Header.Get("Executing"))
		assert.Equal(t, "true", rsp.Header.Get("Executed"))
		assert.Equal(t, "true", rsp.Header.Get("Set-Header1"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header2"))
	})

	t.Run("api/1.0/test/hello3", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8079/api/1.0/testheader/hello3")
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "", apiResponse.Data)
		assert.Equal(t, 200, apiResponse.StatusCode)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, "true", rsp.Header.Get("Executing"))
		assert.Equal(t, "true", rsp.Header.Get("Executed"))
		assert.Equal(t, "true", rsp.Header.Get("Set-Header1"))
		assert.Equal(t, "", rsp.Header.Get("Set-Header2"))
	})
}
