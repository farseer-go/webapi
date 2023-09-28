package middleware

import (
	"github.com/farseer-go/webapi/action"
	"github.com/farseer-go/webapi/context"
	"net/http"
)

// Http HTTP报文响应中间件（默认加载）
type Http struct {
	context.IMiddleware
}

func (receiver *Http) Invoke(httpContext *context.HttpContext) {
	httpContext.Response.SetHttpCode(http.StatusOK)
	httpContext.Response.SetStatusCode(http.StatusOK)

	receiver.IMiddleware.Invoke(httpContext)

	// 说明没有中间件对输出做处理
	if len(httpContext.Response.BodyBytes) == 0 && len(httpContext.Response.Body) > 0 {
		// IActionResult
		if httpContext.IsActionResult() {
			actionResult := httpContext.Response.Body[0].Interface().(action.IResult)
			actionResult.ExecuteResult(httpContext)
		} else {
			// 则转成callResult
			action.NewCallResult().ExecuteResult(httpContext)
		}
	}

	// 输出返回值
	httpContext.Response.W.WriteHeader(httpContext.Response.GetHttpCode())

	// 写入Response流
	if len(httpContext.Response.BodyBytes) > 0 {
		_, _ = httpContext.Response.W.Write(httpContext.Response.BodyBytes)
	}
}
