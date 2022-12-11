package controller

// IActionFilter 过滤器
type IActionFilter interface {
	// OnActionExecuting Action执行前
	OnActionExecuting()
	// OnActionExecuted Action执行后
	OnActionExecuted()
}
