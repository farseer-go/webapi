package webapi

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/middleware"
	"net/http"
	"strings"
)

func Run(params ...string) {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	// 初始化中间件
	middleware.InitMiddleware()

	mux := http.NewServeMux()

	// 将路由表注册到http.HandleFunc
	handleRoute(mux)

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

	if strings.HasPrefix(addr, ":") {
		flog.Infof("Web服务已启动：http://localhost%s/", addr)
	}
	flog.Info(http.ListenAndServe(addr, mux))
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
