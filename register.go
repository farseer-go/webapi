package webapi

import (
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
	"os"
	"strings"
)

func RegisterMiddleware(m middleware.IMiddleware) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	middleware.MiddlewareList.Add(m)
}

// RegisterController 自动注册控制器下的所有Action方法
func RegisterController(c controller.IController) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	controller.Register(defaultApi.area, c)
}

// registerAction 注册单个Api
func registerAction(route Route) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	route.Url = strings.Trim(route.Url, " ")
	route.Url = strings.TrimLeft(route.Url, "/")
	if route.Url == "" {
		flog.Errorf("注册minimalApi失败：%s必须设置值", flog.Colors[eumLogLevel.Error]("route"))
		os.Exit(1)
	}
	minimal.Register(defaultApi.area, route.Method, route.Url, route.Action, route.Params...)
}

// RegisterPOST 注册单个Api
func RegisterPOST(route string, actionFunc any, params ...string) {
	registerAction(Route{
		Url:    route,
		Method: "POST",
		Action: actionFunc,
		Params: params,
	})
}

// RegisterGET 注册单个Api
func RegisterGET(route string, actionFunc any, params ...string) {
	registerAction(Route{
		Url:    route,
		Method: "GET",
		Action: actionFunc,
		Params: params,
	})
}

// RegisterPUT 注册单个Api
func RegisterPUT(route string, actionFunc any, params ...string) {
	registerAction(Route{
		Url:    route,
		Method: "PUT",
		Action: actionFunc,
		Params: params,
	})
}

// RegisterDELETE 注册单个Api
func RegisterDELETE(route string, actionFunc any, params ...string) {
	registerAction(Route{
		Url:    route,
		Method: "DELETE",
		Action: actionFunc,
		Params: params,
	})
}

// RegisterRoutes 批量注册路由
func RegisterRoutes(routes []Route) {
	for i := 0; i < len(routes); i++ {
		registerAction(routes[i])
	}
}
