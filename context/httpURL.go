package context

import (
	"net/http"
	"strings"
)

type HttpURL struct {
	Url         string // 请求地址
	Path        string // 请求地址
	RemoteAddr  string // 客户端IP端口
	Host        string
	Proto       string // http协议
	RequestURI  string
	QueryString string
	Query       map[string]any

	R *http.Request
}

func (r *HttpURL) parseQuery() {
	for k, v := range r.R.URL.Query() {
		key := strings.ToLower(k)
		r.Query[key] = strings.Join(v, "&")
	}
}
