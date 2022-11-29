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
	return receiver.w.Write(content)
}
