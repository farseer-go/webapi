package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
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
	// 场景二：客户端连接后，服务端根据条件多次返回消息
	webapi.RegisterRoutes(webapi.Route{Url: "/ws/api", Method: "WS",
		Action: func(context *websocket.Context[pageSizeRequest]) {
			// 验证头部
			val := context.GetHeader("Token")
			assert.Equal(t, "farseer-go", val)

			for {
				req := context.Receiver()
				context.Send("我收到消息啦：")
				req.PageSize++
				req.PageIndex++
				context.Send(req)
				//context.Close()
			}
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

	t.Run("/ws/api", func(t *testing.T) {
		client, err := ws.NewClient("ws://127.0.0.1:8096/ws/api", 1024)
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
		msg, err := client.Receiver()
		assert.Nil(t, err)
		assert.Equal(t, msg, "我收到消息啦：")

		// 接收服务端的消息
		var request2 pageSizeRequest
		err = client.ReceiverJson(&request2)
		assert.Nil(t, err)
		assert.Equal(t, 201, request2.PageSize)
		assert.Equal(t, 101, request2.PageIndex)

		client.Close()
	})
}
