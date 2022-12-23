package context

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/collections"
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

// BuildActionInValue 根据method映射入参
func (httpContext *HttpContext) BuildActionInValue() []reflect.Value {
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
	return httpContext.Route.mapToParams(mapVal)
}

// application/x-www-form-urlencoded
func (httpContext *HttpContext) formUrlencoded() []reflect.Value {
	// 多参数
	return httpContext.Route.mapToParams(httpContext.Request.Form)
}

// query
func (httpContext *HttpContext) query() []reflect.Value {
	// 多参数
	return httpContext.Route.mapToParams(httpContext.Request.Query)
}

// IsActionResult 是否为ActionResult类型
func (httpContext *HttpContext) IsActionResult() bool {
	return httpContext.Route.ResponseBodyType.Count() == 1 && httpContext.Route.ResponseBodyType.First().String() == "action.IResult"
}
