package context

import "github.com/farseer-go/collections"

type HttpURL struct {
	Url         string // 请求地址
	Path        string // 请求地址
	RemoteAddr  string // 客户端IP端口
	Host        string
	Proto       string // http协议
	RequestURI  string
	QueryString string
	Query       collections.Dictionary[string, string]
}
