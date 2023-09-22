package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	// 测试生成出来的token与head是否一致
	var buildToken string
	configure.SetDefault("WebApi.Jwt.Key", "123456888")
	configure.SetDefault("WebApi.Jwt.KeyType", "HS256")
	configure.SetDefault("WebApi.Jwt.HeaderName", "Auto_test")
	configure.SetDefault("WebApi.Jwt.InvalidStatusCode", 403)
	configure.SetDefault("WebApi.Jwt.InvalidMessage", "您没有权限访问")
	fs.Initialize[webapi.Module]("demo")

	webapi.RegisterRoutes(webapi.Route{Url: "/jwt/build", Action: func() {
		buildToken, _ = webapi.GetHttpContext().Jwt.Build()
	}}.POST())

	webapi.RegisterRoutes(webapi.Route{Url: "/jwt/validate", Action: func() string {
		return "hello"
	}}.UseJwt())

	go webapi.Run(":8090")
	time.Sleep(10 * time.Millisecond)

	t.Run("test jwt build", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8090/jwt/build", "application/json", nil)
		_ = rsp.Body.Close()
		token := rsp.Header.Get(context.HeaderName)
		assert.Equal(t, token, buildToken)
	})

	t.Run("test jwt validate", func(t *testing.T) {
		client := fasthttp.Client{}
		request := fasthttp.AcquireRequest()
		request.SetRequestURI("http://127.0.0.1:8090/jwt/validate")
		request.Header.SetContentType("application/json")
		//request.Header.Set(context.HeaderName, buildToken)
		request.Header.SetMethod("POST")
		response := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(request)
		defer fasthttp.ReleaseResponse(response)
		defer request.SetConnectionClose()
		_ = client.DoTimeout(request, response, 2000*time.Millisecond)

		assert.Equal(t, "hello", string(response.Body()))
		assert.Equal(t, 200, response.StatusCode())
	})
}
