package context

import (
	"net/http"
	"reflect"

	"github.com/farseer-go/fs/snc"
)

type HttpResponse struct {
	W             http.ResponseWriter
	Body          []any  // Action执行的结果（Action返回值）
	BodyBytes     []byte // 自定义输出结果
	httpCode      int    // http响应代码
	statusCode    int    // ApiResponse响应代码
	statusMessage string // ApiResponse响应提示
}

// GetHttpCode 获取响应的HttpCode
func (receiver *HttpResponse) GetHttpCode() int {
	return receiver.httpCode
}

// SetHttpCode 将响应状态写入http流
func (receiver *HttpResponse) SetHttpCode(httpCode int) {
	receiver.httpCode = httpCode
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
	receiver.BodyBytes, _ = snc.Marshal(content)
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

// SetStatusCode 设置StatusCode
func (receiver *HttpResponse) SetStatusCode(statusCode int) {
	receiver.statusCode = statusCode
}

// SetMessage 设计响应提示信息
func (receiver *HttpResponse) SetMessage(statusCode int, statusMessage string) {
	receiver.statusCode = statusCode
	receiver.statusMessage = statusMessage
}

// Reject 拒绝服务
func (receiver *HttpResponse) Reject(httpCode int, content string) {
	receiver.httpCode = httpCode
	receiver.statusCode = httpCode
	receiver.BodyBytes = []byte(content)
}

// GetStatus 获取statusCode、statusMessage
func (receiver *HttpResponse) GetStatus() (int, string) {
	return receiver.statusCode, receiver.statusMessage
}

// SetValues 设置Body值
func (receiver *HttpResponse) SetValues(callValues ...reflect.Value) {
	for _, value := range callValues {
		receiver.Body = append(receiver.Body, value.Interface())
	}
}

// Redirect302 302重定向
func (receiver *HttpResponse) Redirect302(location string) {
	receiver.httpCode = 302
	receiver.W.Header().Set("Location", location)
}

// Redirect301 301重定向
func (receiver *HttpResponse) Redirect301(location string) {
	receiver.httpCode = 301
	receiver.W.Header().Set("Location", location)
}
