package middleware

import "github.com/farseer-go/webapi/context"

// Http HTTP报文响应中间件（默认加载）
type Http struct {
	IMiddleware
}

func (receiver *Http) Invoke(httpContext *context.HttpContext) {
	receiver.IMiddleware.Invoke(httpContext)

	// 输出返回值
	httpContext.Response.WriteCode(httpContext.Response.StatusCode)

	// 有返回值，但没有转成字节
	if len(httpContext.Response.Body) > 0 && len(httpContext.Response.BodyBytes) == 0 {
		// 初始化返回报文
		httpContext.Response.BuildResponse(httpContext.Route)
	}

	// 写入Response流
	if len(httpContext.Response.BodyBytes) > 0 {
		_, _ = httpContext.Response.Write(httpContext.Response.BodyBytes)
	}
}
