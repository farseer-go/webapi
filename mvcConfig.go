package webapi

// 配置
type mvcConfig struct {
	apiPrefix  string // api前缀
	enableCORS bool   // 启用CORS
}

var config = mvcConfig{}
