package context

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/farseer-go/fs/snc"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpRequest struct {
	Body      io.ReadCloser
	BodyBytes []byte // 接收到的字节
	Form      map[string]any
	Query     map[string]any
	Params    []reflect.Value // 转换成Handle函数需要的参数
	R         *http.Request
}

// jsonToMap 将json转成map类型
func (r *HttpRequest) jsonToMap() map[string]any {
	mapVal := make(map[string]any)
	snc.Unmarshal(r.BodyBytes, &mapVal)

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

// msgpackToMap 将 msgpack 二进制转成 map 类型
func (r *HttpRequest) msgpackToMap() map[string]any {
	mapVal := make(map[string]any)

	// 直接将二进制 Body 反序列化到 map 中
	err := msgpack.Unmarshal(r.BodyBytes, &mapVal)
	if err != nil {
		// 如果解包失败，返回空 map 防止后续空指针
		return make(map[string]any)
	}

	// 保持和你 jsonToMap 一致的逻辑：将 Key 转小写
	// 这是为了兼容你后续 FormToParams 的匹配逻辑
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
	if len(r.BodyBytes) > 0 {
		parseQuery, _ := url.ParseQuery(string(r.BodyBytes))
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
