package websocket

import (
	ctx "context"
	"errors"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/context"
	"golang.org/x/net/websocket"
	"net"
	"reflect"
	"time"
)

// Context websocket上下文
type Context[T any] struct {
	Ctx         ctx.Context          // 用于通知应用端是否断开连接
	cancel      ctx.CancelFunc       // 用于通知Ctx，连接已断开
	AutoExit    bool                 // 当断开连接时，自动退出
	HttpContext *context.HttpContext // 上下文
	tType       reflect.Type         // 得到泛型T的类型
	isClose     bool                 // 是否断开连接
}

// ItemType 泛型类型
func (receiver *Context[T]) ItemType() T {
	var t T
	return t
}

// SetContext 收到请求时，设置上下文（webapi使用）
func (receiver *Context[T]) SetContext(httpContext *context.HttpContext) {
	receiver.Ctx, receiver.cancel = ctx.WithCancel(ctx.Background())
	receiver.HttpContext = httpContext
	var t T
	receiver.tType = reflect.TypeOf(t)
	receiver.AutoExit = true
}

// ReceiverMessage 接收消息
func (receiver *Context[T]) ReceiverMessage() string {
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
func (receiver *Context[T]) ReceiverMessageFunc(d time.Duration, f func(message string)) {
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
		f(message)

		// 异步执行函数f
		go func() {
			for {
				select {
				case <-c.Done():
					return
				case <-receiver.Ctx.Done():
					return
				case <-time.Tick(d):
					f(message)
				}
			}
		}()
	}
}

// ReceiverFunc 接收消息。当收到消息后，会执行f()
func (receiver *Context[T]) ReceiverFunc(d time.Duration, f func(message *T)) {
	var c ctx.Context
	var cancel ctx.CancelFunc

	for {
		// 等待消息
		message := receiver.Receiver()
		// 停止上一轮的函数f
		if cancel != nil {
			cancel()
		}
		c, cancel = ctx.WithCancel(ctx.Background())
		f(&message)

		// 异步执行函数f
		go func() {
			for {
				select {
				case <-c.Done():
					return
				case <-receiver.Ctx.Done():
					return
				case <-time.Tick(d):
					f(&message)
				}
			}
		}()
	}
}

// Receiver 接收消息
func (receiver *Context[T]) Receiver() T {
reopen:
	var t T
	switch receiver.tType.Kind() {
	case reflect.Struct:
		err := websocket.JSON.Receive(receiver.HttpContext.WebsocketConn, &t)
		if err != nil {
			receiver.errorIsClose(err)
			flog.Warningf("路由：%s 接收数据时，出现反序列失败：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
			goto reopen
		}
	default:
		var data string
		err := websocket.Message.Receive(receiver.HttpContext.WebsocketConn, &data)
		if err != nil {
			receiver.errorIsClose(err)
			flog.Warningf("路由：%s 接收数据时，出现异常：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
			goto reopen
		}
		return parse.Convert(data, t)
	}
	return t
}

// Send 发送消息，如果msg不是go的基础类型，则会自动序列化成json
func (receiver *Context[T]) Send(msg any) error {
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
func (receiver *Context[T]) Close() {
	_ = receiver.HttpContext.WebsocketConn.Close()
	receiver.cancel()
	receiver.isClose = true
}

// GetHeader 获取头部
func (receiver *Context[T]) GetHeader(key string) string {
	return receiver.HttpContext.Header.GetValue(key)
}

// IsClose 是否已断开连接
func (receiver *Context[T]) IsClose() bool {
	return receiver.isClose
}

// 根据错误信息，判断是否为断开连接导致的
func (receiver *Context[T]) errorIsClose(err error) {
	var opError *net.OpError
	if errors.As(err, &opError) || err.Error() == "EOF" {
		receiver.cancel()
		receiver.isClose = true
		if receiver.AutoExit {
			exception.ThrowWebException(408, "客户端已关闭")
		}
	}
}
