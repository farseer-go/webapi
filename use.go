package webapi

import (
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/middleware"
	"net/http"
)

// UseCors 使用CORS中间件
func UseCors() {
	RegisterMiddleware(&middleware.Cors{})
}

// UseStaticFiles 支持静态目录，在根目录./wwwroot中的文件，直接以静态文件提供服务
func UseStaticFiles() {
	// 需要先依赖模块
	modules.ThrowIfNotLoad(Module{})

	// 默认wwwroot为静态目录
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./wwwroot"))))
}

func UseWebApi() {
	RegisterMiddleware(&middleware.Session{})
	RegisterMiddleware(&middleware.UrlRewriting{})
}

// UseApiResponse 支持ApiResponse结构
func UseApiResponse() {
	RegisterMiddleware(&middleware.ApiResponse{})
}
