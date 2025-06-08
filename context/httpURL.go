package context

import (
	"net/http"
	"strings"

	"github.com/farseer-go/fs/parse"
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
