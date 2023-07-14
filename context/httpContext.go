package context

import (
	"bytes"
	"github.com/farseer-go/collections"
	"net/http"
	"reflect"
	"strings"
)

type HttpContext struct {
	Request          *HttpRequest
	Response         *HttpResponse
	Header           collections.ReadonlyDictionary[string, string]
	Cookie           *HttpCookies
	Session          Session
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
			R:          r,
		},
		Response: &HttpResponse{
			W: w,
		},
		URI: &HttpURL{
			Path:        r.URL.Path,
			RemoteAddr:  r.RemoteAddr,
			Host:        r.Host,
			Proto:       r.Proto,
			RequestURI:  r.RequestURI,
			QueryString: r.URL.RawQuery,
			Query:       collections.NewDictionary[string, string](),
			Url:         "https://" + r.Host + r.RequestURI, // 先默认https，后边在处理
		},
		Method:           r.Method,
		ContentLength:    r.ContentLength,
		Close:            r.Close,
		TransferEncoding: r.TransferEncoding,
		ContentType:      "",
		Route:            &httpRoute,
		Cookie: &HttpCookies{
			w: w,
			r: r,
		},
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
	}

	// header
	header := collections.NewDictionary[string, string]()
	for k, v := range r.Header {
		header.Add(k, strings.Join(v, ";"))
	}
	httpContext.Header = header.ToReadonlyDictionary()
	httpContext.Session, _ = DefaultSession.storage.InitSession("sid", DefaultTime)

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
		return httpContext.Route.JsonToParams(httpContext.Request)
	case "": // GET
		return httpContext.Route.FormToParams(httpContext.Request.Query)
	default: //case "application/x-www-form-urlencoded", "multipart/form-data":
		return httpContext.Route.FormToParams(httpContext.Request.Form)
	}
}

// IsActionResult 是否为ActionResult类型
func (httpContext *HttpContext) IsActionResult() bool {
	return httpContext.Route.ResponseBodyType.Count() == 1 && httpContext.Route.ResponseBodyType.First().String() == "action.IResult"
}
