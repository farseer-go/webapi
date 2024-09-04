package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/utils/ws"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestWebsocket(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	webapi.UseStaticFiles()
	webapi.UseApiResponse()
	webapi.RegisterRoutes(webapi.Route{Url: "/mini/api", Method: "GET", Action: func() any {
		return pageSizeRequest{PageSize: 3, PageIndex: 2}
	}})

	// 场景一：客户端发一次消息，服务端返回一次消息
	webapi.RegisterRoutes(webapi.Route{Url: "/ws/api1", Method: "WS", Params: []string{"context", ""},
		Action: func(context *websocket.Context[pageSizeRequest], manager trace.IManager) {
			// 验证头部
			val := context.GetHeader("Token")
			assert.Equal(t, "farseer-go", val)

			// 验证a、c
			assert.Equal(t, "b", context.HttpContext.Request.Form["a"])
			assert.Equal(t, "d", context.HttpContext.Request.Form["c"])

			req := context.Receiver()
			_ = context.Send("我收到消息啦：")

			req.PageSize++
			req.PageIndex++
			_ = context.Send(req)
		}})

	// 场景二：客户端连接后，客户端每发一次消息，服务端持续返回新的消息
	webapi.RegisterRoutes(webapi.Route{Url: "/ws/api2", Method: "WS", Params: []string{"context", ""},
		Action: func(context *websocket.Context[pageSizeRequest], manager trace.IManager) {
			context.ReceiverMessageFunc(500*time.Millisecond, func(message string) {
				if message == "1" {
					_ = context.Send("hello")
				}

				if message == "2" {
					_ = context.Send("world")
				}
			})
		}})

	go webapi.Run(":8096")
	time.Sleep(100 * time.Millisecond)

	t.Run("mini/api", func(t *testing.T) {
		rsp, _ := http.Get("http://127.0.0.1:8096/mini/api")
		apiResponse := core.NewApiResponseByReader[pageSizeRequest](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 3, apiResponse.Data.PageSize)
		assert.Equal(t, 2, apiResponse.Data.PageIndex)
		assert.Equal(t, 200, rsp.StatusCode)
		assert.Equal(t, 200, apiResponse.StatusCode)
	})

	t.Run("/ws/api1", func(t *testing.T) {
		client, err := ws.NewClient("ws://127.0.0.1:8096/ws/api1?a=b&c=d", 1024)
		assert.Nil(t, err)

		// 设置头部
		client.SetHeader("token", "farseer-go")

		// 连接
		err = client.Connect()
		assert.Nil(t, err)

		// 发消息
		err = client.Send(pageSizeRequest{
			PageSize:  200,
			PageIndex: 100,
		})
		assert.Nil(t, err)

		// 接收服务端的消息
		msg, err := client.ReceiverMessage()
		assert.Nil(t, err)
		assert.Equal(t, msg, "我收到消息啦：")

		// 接收服务端的消息
		var request2 pageSizeRequest
		err = client.Receiver(&request2)
		assert.Nil(t, err)
		assert.Equal(t, 201, request2.PageSize)
		assert.Equal(t, 101, request2.PageIndex)

		time.Sleep(100 * time.Millisecond)
		// 服务端关闭后，尝试继续接收消息
		assert.Panics(t, func() {
			_, _ = client.ReceiverMessage()
		})
	})

	t.Run("/ws/api2", func(t *testing.T) {
		client, err := ws.NewClient("ws://127.0.0.1:8096/ws/api2", 1024)
		assert.Nil(t, err)

		// 连接
		err = client.Connect()
		assert.Nil(t, err)

		// 发送1
		_ = client.Send("1")

		msg, _ := client.ReceiverMessage()
		assert.Equal(t, msg, "hello")

		msg, _ = client.ReceiverMessage()
		assert.Equal(t, msg, "hello")

		_ = client.Send("2")

		msg, _ = client.ReceiverMessage()
		assert.Equal(t, msg, "world")

		msg, _ = client.ReceiverMessage()
		assert.Equal(t, msg, "world")
	})
}
