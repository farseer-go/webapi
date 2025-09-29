package middleware

import (
	"net/http"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/action"
	"github.com/farseer-go/webapi/context"
)

// Http HTTP报文响应中间件（默认加载）
type Http struct {
	context.IMiddleware
}

func (receiver *Http) Invoke(httpContext *context.HttpContext) {
	httpContext.Response.SetHttpCode(http.StatusOK)
	httpContext.Response.SetStatusCode(http.StatusOK)

	// 下一步：exceptionMiddleware
	receiver.IMiddleware.Invoke(httpContext)

	// 说明没有中间件对输出做处理
	if len(httpContext.Response.BodyBytes) == 0 && len(httpContext.Response.Body) > 0 {
		// IActionResult
		if httpContext.IsActionResult() {
			actionResult := httpContext.Response.Body[0].(action.IResult)
			actionResult.ExecuteResult(httpContext)
		} else {
			// 则转成callResult
			action.NewCallResult().ExecuteResult(httpContext)
		}
	}

	// 输出返回值
	httpContext.Response.W.WriteHeader(httpContext.Response.GetHttpCode())

	// 响应header
	rspHeader := collections.NewDictionary[string, string]()
	for k, v := range httpContext.Response.W.Header() {
		rspHeader.Add(k, strings.Join(v, ";"))
	}
	httpContext.ResponseHeader = rspHeader.ToReadonlyDictionary()

	// 写入Response流
	if len(httpContext.Response.BodyBytes) > 0 {
		_, _ = httpContext.Response.W.Write(httpContext.Response.BodyBytes)
	}
}
