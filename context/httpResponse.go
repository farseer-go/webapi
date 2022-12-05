package context

import (
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
