package middleware

import (
	"github.com/farseer-go/linkTrace"
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

	trackContext := linkTrace.NewWebApi(httpContext.URI.Host, httpContext.URI.Url, httpContext.Method, httpContext.ContentType, httpContext.Header, "", httpContext.URI.GetRealIp())
	linkTrace.SetCurTrace(trackContext)

	// 下一步：exceptionMiddleware
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

	linkTrace.GetCurTrace().StatusCode = httpContext.Response.GetHttpCode()
	linkTrace.GetCurTrace().ResponseBody = string(httpContext.Response.BodyBytes)
	// 结束链路追踪
	trackContext.End()

	// 输出返回值
	httpContext.Response.W.WriteHeader(httpContext.Response.GetHttpCode())

	// 写入Response流
	if len(httpContext.Response.BodyBytes) > 0 {
		_, _ = httpContext.Response.W.Write(httpContext.Response.BodyBytes)
	}
}
