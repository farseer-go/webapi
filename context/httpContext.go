package context

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"net/http"
	"reflect"
	"strings"
)

type HttpContext struct {
	HttpRequest      *HttpRequest
	HttpResponse     *HttpResponse
	HttpHeader       collections.Dictionary[string, string]
	HttpURL          *HttpURL
	Method           string
	ContentLength    int64
	Close            bool
	TransferEncoding []string
	ContentType      string
	HttpRoute        *HttpRoute
	Exception        any
}

// NewHttpContext 初始化上下文
func NewHttpContext(httpRoute HttpRoute, w http.ResponseWriter, r *http.Request) HttpContext {
	// Body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	var httpContext = HttpContext{
		HttpRequest: &HttpRequest{
			Body:       r.Body,
			BodyString: buf.String(),
			BodyBytes:  buf.Bytes(),
		},
		HttpResponse: &HttpResponse{
			w: w,
		},
		HttpHeader: collections.NewDictionary[string, string](),
		HttpURL: &HttpURL{
			Path:        r.URL.Path,
			RemoteAddr:  r.RemoteAddr,
			Host:        r.Host,
			Proto:       r.Proto,
			RequestURI:  r.RequestURI,
			QueryString: r.URL.RawQuery,
			Query:       collections.NewDictionary[string, string](),
		},
		Method:           r.Method,
		ContentLength:    r.ContentLength,
		Close:            r.Close,
		TransferEncoding: r.TransferEncoding,
		ContentType:      "",
		HttpRoute:        &httpRoute,
	}

	switch httpContext.Method {
	case "GET":
		_ = r.ParseForm()
		httpContext.HttpRequest.ParseQuery(r.Form)
	default:
		httpContext.HttpRequest.ParseForm()
	}

	// httpURL
	for k, v := range r.URL.Query() {
		httpContext.HttpURL.Query.Add(k, strings.Join(v, ";"))
	}

	if r.TLS == nil {
		httpContext.HttpURL.Url = "http://" + r.Host + r.RequestURI
	} else {
		httpContext.HttpURL.Url = "https://" + r.Host + r.RequestURI
	}

	// header
	for k, v := range r.Header {
		httpContext.HttpHeader.Add(k, strings.Join(v, ";"))
	}

	// ContentType
	for _, contentType := range strings.Split(httpContext.HttpHeader.GetValue("Content-Type"), ";") {
		if strings.Contains(contentType, "/") {
			httpContext.ContentType = contentType
		}
	}
	return httpContext
}

// GetRequestParam 根据method映射入参
func (httpContext *HttpContext) GetRequestParam() []reflect.Value {
	// 没有入参时，忽略request.body
	if httpContext.HttpRoute.RequestParamType.Count() == 0 {
		return []reflect.Value{}
	}

	// application/json
	switch httpContext.ContentType {
	case "application/json":
		return httpContext.contentTypeJson()
	case "": // GET
		return httpContext.query()
	default: //case "application/x-www-form-urlencoded", "multipart/form-data":
		return httpContext.formUrlencoded()
	}
}

// application/json
func (httpContext *HttpContext) contentTypeJson() []reflect.Value {
	// dto
	if httpContext.HttpRoute.RequestParamIsModel {
		firstParamType := httpContext.HttpRoute.RequestParamType.First() // 先取第一个参数
		val := reflect.New(firstParamType).Interface()
		_ = json.Unmarshal(httpContext.HttpRequest.BodyBytes, val)
		return []reflect.Value{reflect.ValueOf(val).Elem()}
	}

	// 多参数
	mapVal := httpContext.HttpRequest.JsonToMap()
	return httpContext.mapToParams(mapVal)
}

// application/x-www-form-urlencoded
func (httpContext *HttpContext) formUrlencoded() []reflect.Value {
	// 多参数
	return httpContext.mapToParams(httpContext.HttpRequest.Form)
}

// query
func (httpContext *HttpContext) query() []reflect.Value {
	// 多参数
	return httpContext.mapToParams(httpContext.HttpRequest.Query)
}

// 将map转成入参值
func (httpContext *HttpContext) mapToParams(mapVal map[string]any) []reflect.Value {
	// dto模式
	if httpContext.HttpRoute.RequestParamIsModel {
		param := httpContext.HttpRoute.RequestParamType.First()
		paramVal := reflect.New(param).Elem()
		for i := 0; i < param.NumField(); i++ {
			field := param.Field(i)
			if !field.IsExported() {
				continue
			}
			key := strings.ToLower(field.Name)
			kv, exists := mapVal[key]
			if exists {
				defVal := paramVal.Field(i).Interface()
				paramVal.FieldByName(field.Name).Set(reflect.ValueOf(parse.Convert(kv, defVal)))
			}
		}
		return []reflect.Value{paramVal}
	}

	// 多参数
	lstParams := make([]reflect.Value, httpContext.HttpRoute.RequestParamType.Count())
	for i := 0; i < httpContext.HttpRoute.RequestParamType.Count(); i++ {
		defVal := reflect.New(httpContext.HttpRoute.RequestParamType.Index(i)).Elem().Interface()
		if httpContext.HttpRoute.ParamNames.Count() > i {
			paramName := strings.ToLower(httpContext.HttpRoute.ParamNames.Index(i))
			paramVal, _ := mapVal[paramName]
			defVal = parse.Convert(paramVal, defVal)
		}
		lstParams[i] = reflect.ValueOf(defVal)
	}
	return lstParams
}

// BuildResponse 初始化返回报文
func (httpContext *HttpContext) BuildResponse() {
	// 没有返回值，则不响应
	if len(httpContext.HttpResponse.Body) == 0 {
		httpContext.HttpResponse.BodyBytes = []byte{}
		httpContext.HttpResponse.BodyString = ""
		return
	}

	// 只有一个返回值
	if len(httpContext.HttpResponse.Body) == 1 {
		responseBody := httpContext.HttpResponse.Body[0].Interface()
		if httpContext.HttpRoute.ResponseBodyIsModel { // dto
			httpContext.HttpResponse.BodyBytes, _ = json.Marshal(responseBody)
			httpContext.HttpResponse.BodyString = string(httpContext.HttpResponse.BodyBytes)
		} else { // 基本类型直接转string
			httpContext.HttpResponse.BodyString = parse.Convert(responseBody, "")
			httpContext.HttpResponse.BodyBytes = []byte(httpContext.HttpResponse.BodyString)
		}
	}

	// 多个返回值，则转成数组Json
	lst := collections.NewListAny()
	for i := 0; i < len(httpContext.HttpResponse.Body); i++ {
		lst.Add(httpContext.HttpResponse.Body[i].Interface())
	}
	httpContext.HttpResponse.BodyBytes, _ = json.Marshal(lst)
	httpContext.HttpResponse.BodyString = string(httpContext.HttpResponse.BodyBytes)
}