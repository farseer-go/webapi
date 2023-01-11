package context

import (
	"net/http"
	"reflect"
)

type HttpResponse struct {
	Body       []reflect.Value
	BodyString string
	BodyBytes  []byte
	W          http.ResponseWriter
	StatusCode int
}

// WriteCode 将响应状态写入http流
func (receiver HttpResponse) WriteCode(statusCode int) {
	receiver.W.WriteHeader(statusCode)
}

// Write 将响应内容写入http流
func (receiver HttpResponse) Write(content []byte) (int, error) {
	//receiver.w.Header().
	return receiver.W.Write(content)
}

// AddHeader 添加头部
func (receiver HttpResponse) AddHeader(key, value string) {
	receiver.W.Header().Add(key, value)
}

// SetHeader 覆盖头部
func (receiver HttpResponse) SetHeader(key, value string) {
	receiver.W.Header().Set(key, value)
}

// DelHeader 删除头部
func (receiver HttpResponse) DelHeader(key string) {
	receiver.W.Header().Del(key)
}
