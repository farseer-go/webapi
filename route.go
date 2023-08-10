package webapi

// Route 路由配置
type Route struct {
	Url     string   // 路由地址
	Method  string   // Method类型（GET|POST|PUT|DELETE）
	Action  any      // Handle
	Params  []string // Handle的入参
	Message string   // api返回的StatusMessage
}
