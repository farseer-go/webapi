package webapi

// IFilter 过滤器
type IFilter interface {
	OnActionExecuting()
}
