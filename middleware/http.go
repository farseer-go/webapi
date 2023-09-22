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

	if httpContext.Response.StatusCode == 0 {
		httpContext.Response.StatusCode = http.StatusOK
	}

	// 输出返回值
	httpContext.Response.WriteCode(httpContext.Response.StatusCode)

	// 写入Response流
	if len(httpContext.Response.BodyBytes) > 0 {
		_, _ = httpContext.Response.W.Write(httpContext.Response.BodyBytes)
	}
}
