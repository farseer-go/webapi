package middleware

import "github.com/farseer-go/webapi/context"

type Cors struct {
	context.IMiddleware
}

func (receiver *Cors) Invoke(httpContext *context.HttpContext) {
	httpContext.Response.AddHeader("Access-Control-Allow-WebHeaders", httpContext.Header.GetValue("Access-Control-Request-WebHeaders"))
	httpContext.Response.AddHeader("Access-Control-Allow-Methods", httpContext.Header.GetValue("Access-Control-Request-Methods"))
	httpContext.Response.AddHeader("Access-Control-Allow-Credentials", "true")
	httpContext.Response.AddHeader("Access-Control-Max-Age", "86400")

	if httpContext.Header.GetValue("Origin") != "" {
		httpContext.Response.AddHeader("Access-Control-Allow-Origin", httpContext.Header.GetValue("Origin"))
	}

	if httpContext.Method == "OPTIONS" {
		httpContext.Response.SetHttpCode(204)
		return
	}
	receiver.IMiddleware.Invoke(httpContext)
}
