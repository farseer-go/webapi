package webapi

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
	"net/http"
	"net/http/pprof"
	"strings"
)

type applicationBuilder struct {
	area           string
	mux            *serveMux
	certFile       string                                // https证书
	keyFile        string                                // https证书 key
	tls            bool                                  // 是否使用https
	MiddlewareList collections.List[context.IMiddleware] // 注册的中间件
}

func NewApplicationBuilder() *applicationBuilder {
	return &applicationBuilder{
		area:           "/",
		mux:            new(serveMux),
		MiddlewareList: collections.NewList[context.IMiddleware](),
	}
}

func (r *applicationBuilder) RegisterMiddleware(m context.IMiddleware) {
	r.MiddlewareList.Add(m)
}

// RegisterPOST 注册单个Api（支持占位符，例如：/{cateId}/{Id}）
func (r *applicationBuilder) RegisterPOST(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "POST", Action: actionFunc, Params: params})
}

// RegisterGET 注册单个Api（支持占位符，例如：/{cateId}/{Id}）
func (r *applicationBuilder) RegisterGET(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "GET", Action: actionFunc, Params: params})
}

// RegisterPUT 注册单个Api（支持占位符，例如：/{cateId}/{Id}）
func (r *applicationBuilder) RegisterPUT(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "PUT", Action: actionFunc, Params: params})
}

// RegisterDELETE 注册单个Api（支持占位符，例如：/{cateId}/{Id}）
func (r *applicationBuilder) RegisterDELETE(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "DELETE", Action: actionFunc, Params: params})
}

// RegisterRoutes 批量注册路由
func (r *applicationBuilder) RegisterRoutes(routes ...Route) {
	for i := 0; i < len(routes); i++ {
		r.registerAction(routes[i])
	}
}

// RegisterController 自动注册控制器下的所有Action方法
func (r *applicationBuilder) RegisterController(c controller.IController) {
	lst := controller.Register(defaultApi.area, c)
	for i := 0; i < lst.Count(); i++ {
		r.mux.HandleRoute(lst.Index(i))
	}
}

// registerAction 注册单个Api
func (r *applicationBuilder) registerAction(route Route) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	route.Url = strings.Trim(route.Url, " ")
	route.Url = strings.TrimLeft(route.Url, "/")
	if route.Url == "" {
		flog.Panicf("注册minimalApi失败：%s必须设置值", flog.Colors[eumLogLevel.Error]("routing"))
	}
	r.mux.HandleRoute(minimal.Register(defaultApi.area, route.Method, route.Url, route.Action, route.Filters, route.Params...))
}

// Area 设置区域
func (r *applicationBuilder) Area(area string, f func()) {
	if !strings.HasPrefix(area, "/") {
		area = "/" + area
	}
	if !strings.HasSuffix(area, "/") {
		area += "/"
	}
	r.area = area
	// 执行注册
	f()
	// 执行完后，恢复区域为"/"
	r.area = "/"
}

// UseCors 使用CORS中间件
func (r *applicationBuilder) UseCors() {
	r.RegisterMiddleware(&middleware.Cors{})
}

// UseStaticFiles 支持静态目录，在根目录./wwwroot中的文件，直接以静态文件提供服务
func (r *applicationBuilder) UseStaticFiles() {
	// 默认wwwroot为静态目录
	r.mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./wwwroot"))))
}

// UsePprof 是否同时开启pprof
func (r *applicationBuilder) UsePprof() {
	r.mux.HandleFunc("/debug/pprof/", pprof.Index)
	r.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// UseSession 开启Session
func (r *applicationBuilder) UseSession() {
	r.RegisterMiddleware(&middleware.Session{})
	go context.ClearSession()
}

func (r *applicationBuilder) UseWebApi() {
	r.RegisterMiddleware(&middleware.UrlRewriting{})
}

// UseApiResponse 支持ApiResponse结构
func (r *applicationBuilder) UseApiResponse() {
	r.RegisterMiddleware(&middleware.ApiResponse{})
}

// UseTLS 使用https
func (r *applicationBuilder) UseTLS(certFile, keyFile string) {
	r.certFile = certFile
	r.keyFile = keyFile
	r.tls = true
}

// Run 运行Web服务
func (r *applicationBuilder) Run(params ...string) {
	// 设置监听地址
	var addr string
	if len(params) > 0 {
		addr = params[0]
	}
	if addr == "" {
		addr = configure.GetString("WebApi.Url")
		if addr == "" {
			addr = ":8888"
		}
	}
	addr = strings.TrimSuffix(addr, "/")

	// 初始化中间件
	r.mux.initMiddleware(r.MiddlewareList)

	if strings.HasPrefix(addr, ":") {
		if r.tls {
			flog.Infof("Web service is started：https://127.0.0.1%s/", addr)
		} else {
			flog.Infof("Web service is started：http://127.0.0.1%s/", addr)
		}
	}

	if r.tls {
		flog.Info(http.ListenAndServeTLS(addr, r.certFile, r.keyFile, r.mux))
	} else {
		flog.Info(http.ListenAndServe(addr, r.mux))
	}
}
