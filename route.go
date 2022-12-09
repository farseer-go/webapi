package webapi

// Route 路由配置
type Route struct {
	Url    string
	Method string
	Action any
	Params []string
}
