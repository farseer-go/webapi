package websocket

import (
	ctx "context"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/context"
	"github.com/timandy/routine"
	"golang.org/x/net/websocket"
	"reflect"
	"time"
)

// Context websocket上下文
type Context[T any] struct {
	*BaseContext
	tType reflect.Type // 得到泛型T的类型
}

// ItemType 泛型类型
func (receiver *Context[T]) ItemType() T {
	var t T
	return t
}

// SetContext 收到请求时，设置上下文（webapi使用）
func (receiver *Context[T]) SetContext(httpContext *context.HttpContext) {
	var t T
	receiver.tType = reflect.TypeOf(t)
	receiver.BaseContext = &BaseContext{
		HttpContext: httpContext,
		AutoExit:    true,
	}
	receiver.Ctx, receiver.cancel = ctx.WithCancel(ctx.Background())
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
		routine.Go(func() {
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
		})
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
