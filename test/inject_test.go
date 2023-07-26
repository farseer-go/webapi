package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/controller"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

type ITestInject interface {
	Call() string
}

type testInject struct {
	value string
}

func (receiver *testInject) Call() string {
	return receiver.value
}

type testInjectController struct {
	controller.BaseController
}

func (r *testInjectController) Hello1(val ITestInject) string {
	return val.Call()
}

// 测试注入
func TestInject(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	// 注册接口实例
	container.Register(func() ITestInject {
		return &testInject{"ok1"}
	})
	container.Register(func() ITestInject {
		return &testInject{"ok2"}
	}, "ins2")

	webapi.RegisterController(&testInjectController{
		BaseController: controller.BaseController{
			Action: map[string]controller.Action{
				"Hello1": {Method: "GET", Params: "ins2"},
			},
		},
	})

	webapi.RegisterGET("/testInjectMini/testMiniApiInject1", func(str string, val ITestInject) string {
		return val.Call() + str
	})

	webapi.RegisterGET("/testInjectMini/testMiniApiInject2", func(str string, val ITestInject) string {
		return val.Call() + str
	}, "str", "ins2")

	go webapi.Run(":8082")
	time.Sleep(100 * time.Millisecond)

	// 测试MVC模式
	t.Run("/testinject/hello1", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8082/testinject/hello1")
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ok2", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	// 测试mini模式
	t.Run("/testInjectMini/testMiniApiInject1", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8082/testInjectMini/testMiniApiInject1?str=baby")
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ok1baby", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	// 测试mini模式
	t.Run("/testInjectMini/testMiniApiInject2", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8082/testInjectMini/testMiniApiInject2?str=oldbaby")
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "ok2oldbaby", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
