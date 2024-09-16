package websocket

import (
	ctx "context"
	"errors"
	"fmt"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"github.com/timandy/routine"
	"golang.org/x/net/websocket"
	"net"
	"time"
)

type BaseContext struct {
	Ctx         ctx.Context          // 用于通知应用端是否断开连接
	cancel      ctx.CancelFunc       // 用于通知Ctx，连接已断开
	AutoExit    bool                 // 当断开连接时，自动退出
	HttpContext *context.HttpContext // 上下文
	isClose     bool                 // 是否断开连接
}

// ReceiverMessage 接收消息
func (receiver *BaseContext) ReceiverMessage() string {
reopen:
	var data string
	err := websocket.Message.Receive(receiver.HttpContext.WebsocketConn, &data)
	if err != nil {
		receiver.errorIsClose(err)
		flog.Warningf("路由：%s 接收数据时，出现异常：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
		goto reopen
	}
	return data
}

// ReceiverMessageFunc 接收消息。当收到消息后，会执行f()
func (receiver *BaseContext) ReceiverMessageFunc(d time.Duration, f func(message string)) {
	var c ctx.Context
	var cancel ctx.CancelFunc

	for {
		// 等待消息
		message := receiver.ReceiverMessage()
		// 停止上一轮的函数f
		if cancel != nil {
			cancel()
		}
		c, cancel = ctx.WithCancel(ctx.Background())

		// 异步执行函数f
		routine.Go(func() {
			for {
				func() {
					var err error
					// 创建链路追踪上下文
					trackContext := container.Resolve[trace.IManager]().EntryWebSocket(receiver.HttpContext.URI.Host, receiver.HttpContext.URI.Url, receiver.HttpContext.ContentType, receiver.HttpContext.Header.ToMap(), receiver.HttpContext.URI.GetRealIp())
					defer func() {
						trackContext.End(err)
					}()

					trackContext.SetBody(message, 0, "")
					exception.Try(func() {
						f(message)
					}).CatchException(func(exp any) {
						err = fmt.Errorf(fmt.Sprint(exp))
					})
				}()

				// 等待下一次执行
				select {
				case <-c.Done():
					return
				case <-receiver.Ctx.Done():
					return
				case <-time.Tick(d):
				}
			}
		})
	}
}

// Send 发送消息，如果msg不是go的基础类型，则会自动序列化成json
func (receiver *BaseContext) Send(msg any) error {
	var err error
	// 基础类型不需要进行序列化
	if fastReflect.PointerOf(msg).Type == fastReflect.GoBasicType {
		err = websocket.Message.Send(receiver.HttpContext.WebsocketConn, msg)
	} else {
		// 其余类型，一律使用json
		err = websocket.JSON.Send(receiver.HttpContext.WebsocketConn, msg)
	}

	if err != nil {
		receiver.errorIsClose(err)
		flog.Warningf("路由：%s 发送数据时失败：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
	}
	return err
}

// Close 关闭连接
func (receiver *BaseContext) Close() {
	_ = receiver.HttpContext.WebsocketConn.Close()
	receiver.cancel()
	receiver.isClose = true
}

// GetHeader 获取头部
func (receiver *BaseContext) GetHeader(key string) string {
	return receiver.HttpContext.Header.GetValue(key)
}

// IsClose 是否已断开连接
func (receiver *BaseContext) IsClose() bool {
	return receiver.isClose
}

// 根据错误信息，判断是否为断开连接导致的
func (receiver *BaseContext) errorIsClose(err error) {
	var opError *net.OpError
	if errors.As(err, &opError) || err.Error() == "EOF" {
		receiver.cancel()
		receiver.isClose = true
		if receiver.AutoExit {
			exception.ThrowWebException(408, "客户端已关闭")
		}
	}
}

// GetParam 获取来自URL的参数
func (receiver *BaseContext) GetParam(key string) any {
	return receiver.HttpContext.Request.Form[key]
}
