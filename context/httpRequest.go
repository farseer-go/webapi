package context

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

type HttpRequest struct {
	Body       io.ReadCloser
	BodyString string
	BodyBytes  []byte
	Form       map[string]any
	Query      map[string]any
}

func (r *HttpRequest) JsonToMap() map[string]any {
	mapVal := make(map[string]any)
	json.Unmarshal(r.BodyBytes, &mapVal)
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

func (r *HttpRequest) ParseForm() {
	r.Form = make(map[string]any)
	formValues := strings.Split(r.BodyString, "&")
	for _, value := range formValues {
		kv := strings.Split(value, "=")
		key := strings.ToLower(kv[0])
		var value any
		if len(kv) > 1 {
			value = kv[1]
		}
		r.Form[key] = value
	}
}

func (r *HttpRequest) ParseQuery(values url.Values) {
	r.Query = make(map[string]any)

	for k, v := range values {
		key := strings.ToLower(k)
		r.Query[key] = strings.Join(v, "&")
	}
}
