package middleware

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
)

// IMiddleware 中间件
type IMiddleware interface {
	Invoke(httpContext *context.HttpContext)
}

// MiddlewareList 已注册的中间件集合
var MiddlewareList collections.List[IMiddleware]
