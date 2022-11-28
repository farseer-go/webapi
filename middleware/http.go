package middleware

import "github.com/farseer-go/webapi/context"

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
