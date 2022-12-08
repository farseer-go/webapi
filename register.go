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

// RegisterAction 注册单个Api
func RegisterAction(method string, route string, actionFunc any, params ...string) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

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
