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
	Request          *HttpRequest
	Response         *HttpResponse
	Header           collections.Dictionary[string, string]
	Route            *HttpRoute
	URI              *HttpURL
	Method           string
	ContentLength    int64
	Close            bool
	TransferEncoding []string
	ContentType      string
	Exception        any
}

// NewHttpContext 初始化上下文
func NewHttpContext(httpRoute HttpRoute, w http.ResponseWriter, r *http.Request) HttpContext {
	// Body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(r.Body)
	var httpContext = HttpContext{
		Request: &HttpRequest{
			Body:       r.Body,
			BodyString: buf.String(),
			BodyBytes:  buf.Bytes(),
		},
		Response: &HttpResponse{
			w: w,
		},
		Header: collections.NewDictionary[string, string](),
		URI: &HttpURL{
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
		Route:            &httpRoute,
	}

	switch httpContext.Method {
	case "GET":
		_ = r.ParseForm()
		httpContext.Request.ParseQuery(r.Form)
	default:
		httpContext.Request.ParseForm()
	}

	// httpURL
	for k, v := range r.URL.Query() {
		httpContext.URI.Query.Add(k, strings.Join(v, ";"))
	}

	if r.TLS == nil {
		httpContext.URI.Url = "http://" + r.Host + r.RequestURI
	} else {
		httpContext.URI.Url = "https://" + r.Host + r.RequestURI
	}

	// header
	for k, v := range r.Header {
		httpContext.Header.Add(k, strings.Join(v, ";"))
	}

	// ContentType
	for _, contentType := range strings.Split(httpContext.Header.GetValue("Content-Type"), ";") {
		if strings.Contains(contentType, "/") {
			httpContext.ContentType = contentType
		}
	}
	return httpContext
}

// GetRequestParam 根据method映射入参
func (httpContext *HttpContext) GetRequestParam() []reflect.Value {
	// 没有入参时，忽略request.body
	if httpContext.Route.RequestParamType.Count() == 0 {
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
	if httpContext.Route.RequestParamIsModel {
		firstParamType := httpContext.Route.RequestParamType.First() // 先取第一个参数
		val := reflect.New(firstParamType).Interface()
		_ = json.Unmarshal(httpContext.Request.BodyBytes, val)
		return []reflect.Value{reflect.ValueOf(val).Elem()}
	}

	// 多参数
	mapVal := httpContext.Request.JsonToMap()
	return httpContext.mapToParams(mapVal)
}

// application/x-www-form-urlencoded
func (httpContext *HttpContext) formUrlencoded() []reflect.Value {
	// 多参数
	return httpContext.mapToParams(httpContext.Request.Form)
}

// query
func (httpContext *HttpContext) query() []reflect.Value {
	// 多参数
	return httpContext.mapToParams(httpContext.Request.Query)
}

// 将map转成入参值
func (httpContext *HttpContext) mapToParams(mapVal map[string]any) []reflect.Value {
	// dto模式
	if httpContext.Route.RequestParamIsModel {
		param := httpContext.Route.RequestParamType.First()
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
	lstParams := make([]reflect.Value, httpContext.Route.RequestParamType.Count())
	for i := 0; i < httpContext.Route.RequestParamType.Count(); i++ {
		defVal := reflect.New(httpContext.Route.RequestParamType.Index(i)).Elem().Interface()
		if httpContext.Route.ParamNames.Count() > i {
			paramName := strings.ToLower(httpContext.Route.ParamNames.Index(i))
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
	if len(httpContext.Response.Body) == 0 {
		httpContext.Response.BodyBytes = []byte{}
		httpContext.Response.BodyString = ""
		return
	}

	// 只有一个返回值
	if len(httpContext.Response.Body) == 1 {
		responseBody := httpContext.Response.Body[0].Interface()
		if httpContext.Route.ResponseBodyIsModel { // dto
			httpContext.Response.BodyBytes, _ = json.Marshal(responseBody)
			httpContext.Response.BodyString = string(httpContext.Response.BodyBytes)
		} else { // 基本类型直接转string
			httpContext.Response.BodyString = parse.Convert(responseBody, "")
			httpContext.Response.BodyBytes = []byte(httpContext.Response.BodyString)
		}
	}

	// 多个返回值，则转成数组Json
	lst := collections.NewListAny()
	for i := 0; i < len(httpContext.Response.Body); i++ {
		lst.Add(httpContext.Response.Body[i].Interface())
	}
	httpContext.Response.BodyBytes, _ = json.Marshal(lst)
	httpContext.Response.BodyString = string(httpContext.Response.BodyBytes)
}
