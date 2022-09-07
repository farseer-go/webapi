package webapi

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/utils/http"
)

// Run webapi.Run() default run on config:FS.
// webapi.Run("localhost")
// webapi.Run(":8089")
// webapi.Run("127.0.0.1:8089")
func Run(params ...string) {
	param := ""
	if len(params) > 0 && params[0] != "" {
		param = params[0]
	}
	if param == "" {
		param = configure.GetString("WebApi.Url")
	}

	param = http.ClearHttpPrefix(param)
	web.BConfig.CopyRequestBody = true
	web.BeeApp.Run(param)
}

// AutoRouter see HttpServer.AutoRouter
func AutoRouter(c web.ControllerInterface) *web.HttpServer {
	return web.BeeApp.AutoRouter(c)
}
