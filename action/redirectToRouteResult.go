package action

import (
	"github.com/farseer-go/webapi/context"
	"net/http"
)

// RedirectToRouteResult 重定向功能
type RedirectToRouteResult struct {
	url string
}

func (receiver RedirectToRouteResult) ExecuteResult(httpContext *context.HttpContext) {
	httpContext.Response.AddHeader("Location", receiver.url)
	httpContext.Response.StatusCode = http.StatusFound
}

// Redirect 重定向
func Redirect(url string) IResult {
	return RedirectToRouteResult{
		url: url,
	}
}
