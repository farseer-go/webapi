package webapi

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
	"net/http"
	"strings"
)

var defaultApi = NewApplicationBuilder()

type applicationBuilder struct {
	area           string
	mux            *http.ServeMux
	certFile       string                                // https证书
	keyFile        string                                // https证书 key
	tls            bool                                  // 是否使用https
	LstRouteTable  collections.List[context.HttpRoute]   // 注册的路由表
	MiddlewareList collections.List[context.IMiddleware] // 注册的中间件
}

func NewApplicationBuilder() *applicationBuilder {
	return &applicationBuilder{
		area:           "/",
		mux:            http.NewServeMux(),
		LstRouteTable:  collections.NewList[context.HttpRoute](),
		MiddlewareList: collections.NewList[context.IMiddleware](),
	}
}

func (r *applicationBuilder) RegisterMiddleware(m context.IMiddleware) {
	r.MiddlewareList.Add(m)
}

// RegisterController 自动注册控制器下的所有Action方法
func (r *applicationBuilder) RegisterController(c controller.IController) {
	lst := controller.Register(defaultApi.area, c)
	r.LstRouteTable.AddRange(lst.AsEnumerable())
}

// registerAction 注册单个Api
func (r *applicationBuilder) registerAction(route Route) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	route.Url = strings.Trim(route.Url, " ")
	route.Url = strings.TrimLeft(route.Url, "/")
	if route.Url == "" {
		panic(flog.Errorf("注册minimalApi失败：%s必须设置值", flog.Colors[eumLogLevel.Error]("routing")))
	}
	r.LstRouteTable.Add(minimal.Register(defaultApi.area, route.Method, route.Url, route.Action, route.Params...))
}

// RegisterPOST 注册单个Api
func (r *applicationBuilder) RegisterPOST(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "POST", Action: actionFunc, Params: params})
}

// RegisterGET 注册单个Api
func (r *applicationBuilder) RegisterGET(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "GET", Action: actionFunc, Params: params})
}

// RegisterPUT 注册单个Api
func (r *applicationBuilder) RegisterPUT(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "PUT", Action: actionFunc, Params: params})
}

// RegisterDELETE 注册单个Api
func (r *applicationBuilder) RegisterDELETE(route string, actionFunc any, params ...string) {
	r.registerAction(Route{Url: route, Method: "DELETE", Action: actionFunc, Params: params})
}

// RegisterRoutes 批量注册路由
func (r *applicationBuilder) RegisterRoutes(routes ...Route) {
	for i := 0; i < len(routes); i++ {
		r.registerAction(routes[i])
	}
}

// 初始化中间件
func (r *applicationBuilder) initMiddleware() {
	middleware.InitMiddleware(r.LstRouteTable, r.MiddlewareList)
}

// Area 设置区域
func (r *applicationBuilder) Area(area string, f func()) {
	if !strings.HasPrefix(area, "/") {
		r.area = "/" + area
	}
	if !strings.HasSuffix(r.area, "/") {
		r.area += "/"
	}

	// 执行注册
	f()
	// 执行完后，恢复区域为"/"
	r.area = "/"
}

// 将路由表注册到http.HandleFunc
func (r *applicationBuilder) handleRoute() {
	// 遍历路由注册表
	for i := 0; i < r.LstRouteTable.Count(); i++ {
		route := r.LstRouteTable.Index(i)
		r.mux.HandleFunc(route.RouteUrl, func(w http.ResponseWriter, r *http.Request) {
			sw := stopwatch.StartNew()
			// 解析报文、组装httpContext
			httpContext := context.NewHttpContext(route, w, r)

			// 执行第一个中间件
			route.HttpMiddleware.Invoke(&httpContext)
			flog.ComponentInfof("webapi", "%s，%s", route.RouteUrl, sw.GetMicrosecondsText())
		})
	}
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

func (r *applicationBuilder) UseWebApi() {
	r.RegisterMiddleware(&middleware.Session{})
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
	r.initMiddleware()

	// 将路由表注册到http.HandleFunc
	r.handleRoute()

	if strings.HasPrefix(addr, ":") {
		if r.tls {
			flog.Infof("Web服务已启动：https://localhost%s/", addr)
		} else {
			flog.Infof("Web服务已启动：http://localhost%s/", addr)
		}
	}

	if r.tls {
		flog.Info(http.ListenAndServeTLS(addr, r.certFile, r.keyFile, r.mux))
	} else {
		flog.Info(http.ListenAndServe(addr, r.mux))
	}
}
