package webapi

import (
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/filter"
)

// Route 路由配置
type Route struct {
	Method  string            // Method类型（GET|POST|PUT|DELETE）
	Url     string            // 路由地址
	Action  any               // Handle
	Message string            // api返回的StatusMessage
	Filters []context.IFilter // 过滤器（对单个路由的执行单元）
	Params  []string          // Handle的入参
}

// UseJwt 使用Jwt
func (receiver Route) UseJwt() Route {
	receiver.Filters = append(receiver.Filters, &filter.JwtFilter{})
	return receiver
}

// POST 使用POST
func (receiver Route) POST() Route {
	receiver.Method = "POST"
	return receiver
}

// GET 使用GET
func (receiver Route) GET() Route {
	receiver.Method = "GET"
	return receiver
}
