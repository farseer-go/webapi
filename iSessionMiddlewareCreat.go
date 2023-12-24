package webapi

import "github.com/farseer-go/webapi/context"

type ISessionMiddlewareCreat interface {
	// Create 创建Session中间件
	Create() context.IMiddleware
}
