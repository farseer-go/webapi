package webapi

import "github.com/beego/beego/v2/server/web"

// RegisterApi 注册API
func RegisterApi(c web.ControllerInterface) {
	if config.apiPrefix == "" {
		web.BeeApp.AutoRouter(c)
	} else {
		web.BeeApp.AutoPrefix(config.apiPrefix, c)
	}
}

// SetApiPrefix 设置api前缀
func SetApiPrefix(prefix string) {
	config.apiPrefix = prefix
}
