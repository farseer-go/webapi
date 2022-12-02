package context

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"net/http"
	"reflect"
)

type HttpResponse struct {
	Body       []reflect.Value
	BodyString string
	BodyBytes  []byte
	w          http.ResponseWriter
	StatusCode int
}

// WriteCode 将响应状态写入http流
func (receiver HttpResponse) WriteCode(statusCode int) {
	receiver.w.WriteHeader(statusCode)
}

// Write 将响应内容写入http流
func (receiver HttpResponse) Write(content []byte) (int, error) {
	//receiver.w.Header().
	return receiver.w.Write(content)
}

// AddHeader 添加头部
func (receiver HttpResponse) AddHeader(key, value string) {
	receiver.w.Header().Add(key, value)
}

// SetHeader 覆盖头部
func (receiver HttpResponse) SetHeader(key, value string) {
	receiver.w.Header().Set(key, value)
}

// DelHeader 删除头部
func (receiver HttpResponse) DelHeader(key string) {
	receiver.w.Header().Del(key)
}

// 初始化返回报文
func (receiver *HttpResponse) BuildResponse(route *HttpRoute) {
	// 没有返回值，则不响应
	if len(receiver.Body) == 0 {
		receiver.BodyBytes = []byte{}
		receiver.BodyString = ""
		return
	}

	// 只有一个返回值
	if len(receiver.Body) == 1 {
		responseBody := receiver.Body[0].Interface()
		if route.ResponseBodyIsModel { // dto
			receiver.BodyBytes, _ = json.Marshal(responseBody)
			receiver.BodyString = string(receiver.BodyBytes)
		} else { // 基本类型直接转string
			receiver.BodyString = parse.Convert(responseBody, "")
			receiver.BodyBytes = []byte(receiver.BodyString)
		}
		return
	}

	// 多个返回值，则转成数组Json
	lst := collections.NewListAny()
	for i := 0; i < len(receiver.Body); i++ {
		lst.Add(receiver.Body[i].Interface())
	}
	receiver.BodyBytes, _ = json.Marshal(lst)
	receiver.BodyString = string(receiver.BodyBytes)
}
