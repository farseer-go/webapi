package webapi

// Route 路由配置
type Route struct {
	Method  string   // Method类型（GET|POST|PUT|DELETE）
	Url     string   // 路由地址
	Action  any      // Handle
	Message string   // api返回的StatusMessage
	Params  []string // Handle的入参
}
