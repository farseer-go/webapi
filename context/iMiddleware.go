package context

// IMiddleware 中间件
type IMiddleware interface {
	Invoke(httpContext *HttpContext)
}
