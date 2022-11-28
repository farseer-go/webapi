package minimal

import (
	"github.com/farseer-go/webapi/context"
	"reflect"
)

type MinimalMiddleware struct {
}

func (receiver MinimalMiddleware) Invoke(httpContext *context.HttpContext) {
	// 入参
	params := httpContext.GetRequestParam()
	// 调用action
	returnVals := reflect.ValueOf(httpContext.HttpRoute.Action).Call(params)
	// 初始化返回报文
	httpContext.InitResponse(returnVals)
}
