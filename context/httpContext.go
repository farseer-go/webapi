package context

import (
	"github.com/farseer-go/collections"
	"net/http"
	"reflect"
	"strings"
)

type HttpContext struct {
	Request          *HttpRequest                                   // Request
	Response         *HttpResponse                                  // Response
	Header           collections.ReadonlyDictionary[string, string] // 头部信息
	Cookie           *HttpCookies                                   // Cookies信息
	Session          IHttpSession                                   // Session信息
	Route            *HttpRoute                                     // 路由信息
	URI              *HttpURL                                       // URL信息
	Data             *HttpData                                      // 用于传递值
	Method           string                                         // 客户端提交时的Method
	ContentLength    int64                                          // 客户端提交时的内容长度
	ContentType      string                                         // 客户端提交时的内容类型
	Exception        any                                            // 是否发生异常
	Jwt              *HttpJwt                                       // jwt验证
	Close            bool
	TransferEncoding []string
}

// NewHttpContext 初始化上下文
func NewHttpContext(httpRoute *HttpRoute, w http.ResponseWriter, r *http.Request) *HttpContext {
	var httpContext = HttpContext{
		Request: &HttpRequest{
			Body:  r.Body,
			R:     r,
			Form:  make(map[string]any),
			Query: make(map[string]any),
		},
		Response: &HttpResponse{
			W:             w,
			statusMessage: "成功",
		},
		URI: &HttpURL{
			Path:            r.URL.Path,
			RemoteAddr:      r.RemoteAddr,
			X_Forwarded_For: r.Header.Get("X-Forwarded-For"),
			X_Real_Ip:       r.Header.Get("X-Real-Ip"),
			Host:            r.Host,
			Proto:           r.Proto,
			RequestURI:      r.RequestURI,
			QueryString:     r.URL.RawQuery,
			Query:           make(map[string]any),
			Url:             "https://" + r.Host + r.RequestURI, // 先默认https，后边在处理
			R:               r,
		},
		Data:             &HttpData{value: collections.NewDictionary[string, any]()},
		Method:           r.Method,
		ContentLength:    r.ContentLength,
		Close:            r.Close,
		TransferEncoding: r.TransferEncoding,
		ContentType:      "",
		Route:            httpRoute,
		Cookie:           initCookies(w, r),
		Jwt: &HttpJwt{
			w: w,
			r: r,
		},
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

	// WebContentType
	for _, contentType := range strings.Split(httpContext.Header.GetValue("Content-Type"), ";") {
		if strings.Contains(contentType, "/") {
			httpContext.ContentType = contentType
		}
	}
	return &httpContext
}

// ParseParams 转换成Handle函数需要的参数
func (httpContext *HttpContext) ParseParams() []reflect.Value {
	// 没有入参时，忽略request.body
	if httpContext.Route.RequestParamType.Count() == 0 {
		return []reflect.Value{}
	}

	if httpContext.Method == "GET" {
		return httpContext.Route.FormToParams(httpContext.Request.Query)
	}

	// application/json
	switch httpContext.ContentType {
	case "application/json":
		return httpContext.Route.JsonToParams(httpContext.Request)
	default: //case "application/x-www-form-urlencoded", "multipart/form-data":
		return httpContext.Route.FormToParams(httpContext.Request.Query) // Query比Form有更齐全的值，所以不用Form
	}
}

// IsActionResult 是否为ActionResult类型
func (httpContext *HttpContext) IsActionResult() bool {
	return httpContext.Route.ResponseBodyType.Count() == 1 && httpContext.Route.ResponseBodyType.First().String() == "action.IResult"
}
