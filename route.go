package webapi

import (
	"github.com/farseer-go/webapi/context"
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
