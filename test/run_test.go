package test

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")

	webapi.Area("/api/1.0/", func() {

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

	webapi.UseApiResponse()
	go webapi.Run(":8888")
	time.Sleep(100 * time.Millisecond)

	t.Run("hello1", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		request := bytes.NewReader(marshal)
		rsp, _ := http.Post("http://127.0.0.1:8888/api/1.0/mini/hello1", "application/json", request)
		assert.Equal(t, 200, rsp.StatusCode)
		defer rsp.Body.Close()
		rspByte, _ := io.ReadAll(rsp.Body)

		var apiResponse core.ApiResponseString
		_ = json.Unmarshal(rspByte, &apiResponse)
		assert.Equal(t, Hello1(sizeRequest), apiResponse.Data)
	})
}
