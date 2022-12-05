package action

import "github.com/farseer-go/webapi/context"

// RedirectToRouteResult 重定向功能
type RedirectToRouteResult struct {
	url string
}

func (receiver RedirectToRouteResult) ExecuteResult(httpContext *context.HttpContext) {
	httpContext.Response.Redirect(receiver.url)
}

// Redirect 重定向
func Redirect(url string) IResult {
	return RedirectToRouteResult{
		url: url,
	}
}
