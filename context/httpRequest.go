package context

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type HttpRequest struct {
	Body       io.ReadCloser
	BodyString string
	BodyBytes  []byte
	Form       map[string]any
	Query      map[string]any

	R *http.Request
}

// jsonToMap 将json转成map类型
func (r *HttpRequest) jsonToMap() map[string]any {
	mapVal := make(map[string]any)
	_ = json.Unmarshal(r.BodyBytes, &mapVal)
	// 将Key转小写
	for k, v := range mapVal {
		kLower := strings.ToLower(k)
		if k != kLower {
			delete(mapVal, k)
			mapVal[kLower] = v
		}
	}
	return mapVal
}

// 解析来自form的值
func (r *HttpRequest) parseForm() {
	for k, v := range r.R.Form {
		key := strings.ToLower(k)
		r.Form[key] = strings.Join(v, "&")
		r.Query[key] = strings.Join(v, "&")
	}
}

// 解析来自url的值
func (r *HttpRequest) parseQuery() {
	for k, v := range r.R.URL.Query() {
		key := strings.ToLower(k)
		r.Query[key] = strings.Join(v, "&")
	}
}
