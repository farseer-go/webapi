package context

import (
	"net/http"
	"strings"
)

type HttpURL struct {
	Url             string // 请求地址 http(s)://xxx.xxx.xxx/xxx
	Path            string // 请求地址
	RemoteAddr      string // 客户端IP端口
	X_Forwarded_For string // 客户端IP端口
	X_Real_Ip       string // 客户端IP端口
	Host            string // 请求的Host主机头
	Proto           string // http协议
	RequestURI      string
	QueryString     string
	Query           map[string]any
	R               *http.Request
}

//func (receiver *HttpURL) ParseQuery() {
//	for k, v := range receiver.R.URL.Query() {
//		key := strings.ToLower(k)
//		receiver.Query[key] = strings.Join(v, "&")
//	}
//}

// GetRealIp 获取真实IP
func (receiver *HttpURL) GetRealIp() string {
	ip := receiver.X_Real_Ip
	if ip == "" {
		ip = strings.Split(receiver.X_Forwarded_For, ",")[0]
	}
	if ip == "" {
		ip = receiver.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}
