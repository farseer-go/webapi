package websocket

import (
	ctx "context"

	"fmt"
	"reflect"
	"time"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
	"github.com/timandy/routine"
	"golang.org/x/net/websocket"
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
		messageStr, messageData := receiver.receiver()
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
					trackContext := container.Resolve[trace.IManager]().EntryWebSocket(receiver.HttpContext.URI.Host, receiver.HttpContext.URI.Url, receiver.HttpContext.Header.ToMap(), receiver.HttpContext.URI.GetRealIp())
					defer func() {
						container.Resolve[trace.IManager]().Push(trackContext, err)
					}()

					trackContext.SetBody(messageStr, 0, "")
					exception.Try(func() {
						f(&messageData)
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

// ForReceiverFunc 持续接收消息然后执行f()，然后再接收
func (receiver *Context[T]) ForReceiverFunc(f func(message *T)) {
	// 执行函数f
	for {
		// 等待消息
		messageStr, messageData := receiver.receiver()
		var err error
		// 创建链路追踪上下文
		trackContext := container.Resolve[trace.IManager]().EntryWebSocket(receiver.HttpContext.URI.Host, receiver.HttpContext.URI.Url, receiver.HttpContext.Header.ToMap(), receiver.HttpContext.URI.GetRealIp())
		trackContext.SetBody(messageStr, 0, "")
		exception.Try(func() {
			f(&messageData)
		}).CatchException(func(exp any) {
			err = fmt.Errorf(fmt.Sprint(exp))
		})
		container.Resolve[trace.IManager]().Push(trackContext, err)
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

// Receiver 接收消息
func (receiver *Context[T]) receiver() (string, T) {
reopen:
	var message string
	if err := websocket.Message.Receive(receiver.HttpContext.WebsocketConn, &message); err != nil {
		receiver.errorIsClose(err)
		flog.Warningf("路由：%s 接收数据时，出现异常：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
		goto reopen
	}

	// 序列化
	var t T
	if err := snc.Unmarshal([]byte(message), &t); err != nil {
		receiver.errorIsClose(err)
		flog.Warningf("路由：%s 接收数据时，出现反序列失败：%s", receiver.HttpContext.Route.RouteUrl, err.Error())
		goto reopen
	}
	return message, parse.Convert(message, t)
}
