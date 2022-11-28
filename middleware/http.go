package middleware

import "github.com/farseer-go/webapi/context"

// HTTP报文响应中间件（默认加载）
type Http struct {
	IMiddleware
}

func (receiver *Http) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)

	// 输出返回值
	httpContext.HttpResponse.WriteCode(httpContext.HttpResponse.StatusCode)

	if httpContext.HttpResponse.BodyBytes != nil {
		_, _ = httpContext.HttpResponse.Write(httpContext.HttpResponse.BodyBytes)
	}
}
