package context

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type HttpResponse struct {
	W             http.ResponseWriter
	Body          []reflect.Value // Action执行的结果（Action返回值）
	BodyBytes     []byte          // 自定义输出结果
	StatusCode    int             // 响应代码
	StatusMessage string          // 响应提示
}

// WriteCode 将响应状态写入http流
func (receiver *HttpResponse) WriteCode(statusCode int) {
	receiver.W.WriteHeader(statusCode)
}

// Write 将响应内容写入http流
func (receiver *HttpResponse) Write(content []byte) {
	receiver.BodyBytes = content
}

// WriteString 将响应内容写入http流
func (receiver *HttpResponse) WriteString(content string) {
	receiver.BodyBytes = []byte(content)
}

// WriteJson 将响应内容转成json后写入http流
func (receiver *HttpResponse) WriteJson(content any) {
	receiver.BodyBytes, _ = json.Marshal(content)
	receiver.W.Header().Set("Content-Type", "application/json")
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

// Error405 405 Method不被允许访问
func (receiver *HttpResponse) Error405() {
	receiver.StatusCode = http.StatusMethodNotAllowed
	receiver.BodyBytes = []byte("405 Method NotAllowed")
}
