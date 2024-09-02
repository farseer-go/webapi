package websocket

import (
	"encoding/json"
	"errors"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/context"
	"golang.org/x/net/websocket"
	"net"
	"reflect"
)

// Context websocket上下文
type Context[T any] struct {
	httpContext *context.HttpContext
	tType       reflect.Type
}

// ItemType 泛型类型
func (receiver *Context[T]) ItemType() T {
	var t T
	return t
}

// SetContext 收到请求时，设置上下文（webapi使用）
func (receiver *Context[T]) SetContext(httpContext *context.HttpContext) {
	receiver.httpContext = httpContext
	var t T
	receiver.tType = reflect.TypeOf(t)
}

// Receiver 接收消息
func (receiver *Context[T]) Receiver() T {
reopen:
	var t T
	switch receiver.tType.Kind() {
	case reflect.Struct:
		err := websocket.JSON.Receive(receiver.httpContext.WebsocketConn, &t)
		if err != nil {
			var opError *net.OpError
			if errors.As(err, &opError) {
				exception.ThrowWebException(408, "客户端已关闭")
			}

			flog.Warningf("路由：%s 接收数据时，出现反序列失败：%s", receiver.httpContext.Route.RouteUrl, err.Error())
			goto reopen
		}
	default:
		var data string
		err := websocket.Message.Receive(receiver.httpContext.WebsocketConn, &data)
		if err != nil {
			flog.Warningf("路由：%s 接收数据时，出现异常：%s", receiver.httpContext.Route.RouteUrl, err.Error())
			goto reopen
		}
		return parse.Convert(data, t)
	}
	return t
}

// Send 发送消息，如果msg不是go的基础类型，则会自动序列化成json
func (receiver *Context[T]) Send(msg any) {
	switch fastReflect.PointerOf(msg).Type {
	case fastReflect.GoBasicType:
		_ = websocket.Message.Send(receiver.httpContext.WebsocketConn, msg)
	default:
		marshalBytes, err := json.Marshal(msg)
		if err != nil {
			flog.Warningf("路由：%s 发送数据时，出现反序列失败：%s", receiver.httpContext.Route.RouteUrl, err.Error())
		}
		_, _ = receiver.httpContext.WebsocketConn.Write(marshalBytes)
	}
}

// Close 关闭连接
func (receiver *Context[T]) Close() {
	_ = receiver.httpContext.WebsocketConn.Close()
	exception.ThrowWebException(408, "服务端关闭")
}
