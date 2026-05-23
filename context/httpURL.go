package context

import (
	"net/http"
	"strings"

	"github.com/farseer-go/fs/parse"
)

type HttpURL struct {
	Url              string // 请求地址 http(s)://xxx.xxx.xxx/xxx
	Path             string // 请求地址
	RemoteAddr       string // 客户端IP端口
	X_Forwarded_For  string // 客户端IP
	X_Real_Ip        string // 客户端IP
	Cf_Connecting_Ip string // CF提供的真实IP
	Host             string // 请求的Host主机头
	Proto            string // http协议
	RequestURI       string
	QueryString      string
	Query            map[string]any
	R                *http.Request
}

//func (receiver *HttpURL) ParseQuery() {
//	for k, v := range receiver.R.URL.Query() {
//		key := strings.ToLower(k)
//		receiver.Query[key] = strings.Join(v, "&")
//	}
//}

// GetRealIp 获取真实IP
func (receiver *HttpURL) GetRealIp() string {
	ips := []string{receiver.Cf_Connecting_Ip, receiver.X_Real_Ip, strings.Split(receiver.X_Forwarded_For, ",")[0], receiver.RemoteAddr}
	for _, ip := range ips {
		if ip = strings.TrimSpace(ip); ip != "" {
			return strings.SplitN(ip, ":", 2)[0]
		}
	}
	return "" // 极端情况下（如全空）的最终安全兜底
}

// GetRealIpPort 获取真实IP、Port
func (receiver *HttpURL) GetRealIpPort() (string, int) {
	ips := []string{receiver.X_Real_Ip, strings.Split(receiver.X_Forwarded_For, ",")[0], receiver.RemoteAddr}
	for _, ip := range ips {
		if ipPorts := strings.Split(ip, ":"); len(ipPorts) == 2 {
			return ipPorts[0], parse.ToInt(ipPorts[1])
		}
	}

	// 如果没有找到IP和端口，则返回RemoteAddr的IP部分和0端口
	return strings.Split(receiver.RemoteAddr, ":")[0], 0
}
