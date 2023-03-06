package webapi

import (
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
)

func RegisterMiddleware(m context.IMiddleware) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})
	defaultApi.RegisterMiddleware(m)
}

// RegisterController 自动注册控制器下的所有Action方法
func RegisterController(c controller.IController) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})
	defaultApi.RegisterController(c)
}

// RegisterPOST 注册单个Api
func RegisterPOST(route string, actionFunc any, params ...string) {
	defaultApi.RegisterPOST(route, actionFunc, params...)
}

// RegisterGET 注册单个Api
func RegisterGET(route string, actionFunc any, params ...string) {
	defaultApi.RegisterGET(route, actionFunc, params...)
}

// RegisterPUT 注册单个Api
func RegisterPUT(route string, actionFunc any, params ...string) {
	defaultApi.RegisterPUT(route, actionFunc, params...)
}

// RegisterDELETE 注册单个Api
func RegisterDELETE(route string, actionFunc any, params ...string) {
	defaultApi.RegisterDELETE(route, actionFunc, params...)
}

// RegisterRoutes 批量注册路由
func RegisterRoutes(routes ...Route) {
	defaultApi.RegisterRoutes(routes...)
}

// Area 设置区域
func Area(area string, f func()) {
	defaultApi.Area(area, f)
}

// UseCors 使用CORS中间件
func UseCors() {
	defaultApi.UseCors()
}

// UseStaticFiles 支持静态目录，在根目录./wwwroot中的文件，直接以静态文件提供服务
func UseStaticFiles() {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	defaultApi.UseStaticFiles()
}

// UsePprof 是否同时开启pprof
func UsePprof() {
	defaultApi.UsePprof()
}

func UseWebApi() {
	defaultApi.UseWebApi()
}

// UseApiResponse 让所有的返回值，包含在core.ApiResponse中
func UseApiResponse() {
	defaultApi.UseApiResponse()
}

// UseTLS 使用https
func UseTLS(certFile, keyFile string) {
	defaultApi.UseTLS(certFile, keyFile)
}

// Run 运行Web服务
func Run(params ...string) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})
	defaultApi.Run(params...)
}
