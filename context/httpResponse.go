package context

import (
	"net/http"
	"reflect"
)

type HttpResponse struct {
	Body          []reflect.Value
	BodyString    string
	BodyBytes     []byte
	W             http.ResponseWriter
	StatusCode    int    // 响应代码
	StatusMessage string // 响应提示
}

// WriteCode 将响应状态写入http流
func (receiver *HttpResponse) WriteCode(statusCode int) {
	receiver.W.WriteHeader(statusCode)
}

// Write 将响应内容写入http流
func (receiver *HttpResponse) Write(content []byte) (int, error) {
	return receiver.W.Write(content)
}

// AddHeader 添加头部
func (receiver *HttpResponse) AddHeader(key, value string) {
	receiver.W.Header().Add(key, value)
}

// SetHeader 覆盖头部
func (receiver *HttpResponse) SetHeader(key, value string) {
	receiver.W.Header().Set(key, value)
}

// DelHeader 删除头部
func (receiver *HttpResponse) DelHeader(key string) {
	receiver.W.Header().Del(key)
}

// SetMessage 设计响应提示信息
func (receiver *HttpResponse) SetMessage(statusMessage string) {
	receiver.StatusMessage = statusMessage
}
