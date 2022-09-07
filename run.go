package webapi

import "github.com/beego/beego/v2/server/web"

// Run webapi.Run() default run on HttpPort
// webapi.Run("localhost")
// webapi.Run(":8089")
// webapi.Run("127.0.0.1:8089")
func Run(params ...string) {
	if len(params) > 0 && params[0] != "" {
		web.BeeApp.Run(params[0])
	}
	web.BeeApp.Run("")
}

// AutoRouter see HttpServer.AutoRouter
func AutoRouter(c web.ControllerInterface) *web.HttpServer {
	return web.BeeApp.AutoRouter(c)
}
