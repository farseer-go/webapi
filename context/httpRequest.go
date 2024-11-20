package context

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type HttpRequest struct {
	Body       io.ReadCloser
	BodyString string
	BodyBytes  []byte
	Form       map[string]any
	Query      map[string]any
	Params     []reflect.Value // 转换成Handle函数需要的参数
	R          *http.Request
}

// jsonToMap 将json转成map类型
func (r *HttpRequest) jsonToMap() map[string]any {
	mapVal := make(map[string]any)
	//_ = sonic.Unmarshal(r.BodyBytes, &mapVal)
	d := sonic.NewDecoder(bytes.NewReader(r.BodyBytes))
	d.UseNumber()
	_ = d.Decode(&mapVal)

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

// ParseForm 解析来自form的值
func (r *HttpRequest) ParseForm() {
	for k, v := range r.R.Form {
		key := strings.ToLower(k)
		r.Form[key] = strings.Join(v, "&")
		r.Query[key] = strings.Join(v, "&")
	}

	// multipart/form-data提交的数据在Body中
	if r.BodyString != "" {
		parseQuery, _ := url.ParseQuery(r.BodyString)
		for key, value := range parseQuery {
			key = strings.ToLower(key)
			if len(value) > 0 {
				r.Form[key] = strings.Join(value, ",")
				r.Query[key] = strings.Join(value, ",")
			}
		}
		//formValues := strings.Split(r.BodyString, "&")
		//for _, value := range formValues {
		//	kv := strings.Split(value, "=")
		//	if len(kv) > 1 {
		//		key := strings.ToLower(kv[0])
		//		r.Form[key] = kv[1]
		//		r.Query[key] = kv[1]
		//	}
		//}
	}
}

// ParseQuery 解析来自url的值
func (r *HttpRequest) ParseQuery() {
	for k, v := range r.R.URL.Query() {
		key := strings.ToLower(k)
		r.Query[key] = strings.Join(v, "&")
	}
}
