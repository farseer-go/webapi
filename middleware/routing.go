package middleware

import (
	"bytes"
	"github.com/farseer-go/webapi/context"
)

type routing struct {
	context.IMiddleware
}

func (receiver *routing) Invoke(httpContext *context.HttpContext) {
	// 检查method
	if httpContext.Method != "OPTIONS" && !httpContext.Route.Method.Contains(httpContext.Method) {
		// 响应码
		httpContext.Response.Reject(405, "405 Method NotAllowed")
		return
	}

	// 解析Body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(httpContext.Request.R.Body)
	httpContext.Request.BodyString = buf.String()
	httpContext.Request.BodyBytes = buf.Bytes()

	// 解析请求的参数
	//httpContext.URI.ParseQuery()
	httpContext.Request.ParseQuery()
	httpContext.Request.ParseForm()

	// 转换成Handle函数需要的参数
	httpContext.Request.Params = httpContext.ParseParams()

	receiver.IMiddleware.Invoke(httpContext)
}
