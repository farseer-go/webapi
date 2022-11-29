package webapi

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
	"net/http"
	"os"
	"strings"
)

func Run(params ...string) {
	// 初始化中间件
	middleware.Init()

	// 处理路由
	handleRoute()

	var addr string
	if len(params) > 0 && params[0] != "" {
		addr = params[0]
	}
	if addr == "" {
		addr = configure.GetString("WebApi.Url")
	}

	if addr == "" {
		addr = ":8888"
	}

	if strings.HasPrefix(addr, ":") {
		flog.Infof("Web服务已启动：http://localhost%s/", addr)
	}
	flog.Info(http.ListenAndServe(addr, nil))
}

func RegisterMiddleware(m middleware.IMiddleware) {
	middleware.MiddlewareList.Add(m)
}

// Area 设置区域
func Area(area string, f func()) {
	if !strings.HasPrefix(area, "/") {
		area = "/" + area
	}
	if !strings.HasSuffix(area, "/") {
		area += "/"
	}
	defaultApi.area = area
	// 执行注册
	f()
	// 执行完后，恢复区域为"/"
	defaultApi.area = "/"
}

// RegisterController 自动注册控制器下的所有Action方法
func RegisterController(c controller.IController) {
	controller.Register(defaultApi.area, c)
}

// RegisterAction 注册单个Api
func RegisterAction(method string, route string, actionFunc any, params ...string) {
	route = strings.Trim(route, " ")
	route = strings.TrimLeft(route, "/")
	if route == "" {
		flog.Errorf("注册minimalApi失败：%s必须设置值", flog.Colors[eumLogLevel.Error]("route"))
		os.Exit(1)
	}
	minimal.Register(defaultApi.area, method, route, actionFunc, params...)
}

// RegisterPOST 注册单个Api
func RegisterPOST(route string, actionFunc any, params ...string) {
	RegisterAction("POST", route, actionFunc, params...)
}

// RegisterGET 注册单个Api
func RegisterGET(route string, actionFunc any, params ...string) {
	RegisterAction("GET", route, actionFunc, params...)
}

// RegisterPUT 注册单个Api
func RegisterPUT(route string, actionFunc any, params ...string) {
	RegisterAction("PUT", route, actionFunc, params...)
}

// RegisterDELETE 注册单个Api
func RegisterDELETE(route string, actionFunc any, params ...string) {
	RegisterAction("DELETE", route, actionFunc, params...)
}

// UseCors 使用CORS中间件
func UseCors() {
	RegisterMiddleware(&middleware.Cors{})
}

// UseStaticFiles 支持静态目录，在根目录./wwwroot中的文件，直接以静态文件提供服务
func UseStaticFiles() {
	// 默认wwwroot为静态目录
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./wwwroot"))))
}

func UseWebApi() {
	RegisterMiddleware(&middleware.Routing{})
	RegisterMiddleware(&middleware.Session{})
	RegisterMiddleware(&middleware.UrlRewriting{})
}

// UseApiResponse 支持ApiResponse结构
func UseApiResponse() {
	RegisterMiddleware(&middleware.ApiResponse{})
}
