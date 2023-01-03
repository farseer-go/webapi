package action

import (
	"encoding/json"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/context"
)

// callResult 默认调用Action结果
type callResult struct {
}

func NewCallResult() callResult {
	return callResult{}
}

func (receiver callResult) ExecuteResult(httpContext *context.HttpContext) {
	// 只有一个返回值
	if len(httpContext.Response.Body) == 1 {
		responseBody := httpContext.Response.Body[0].Interface()
		// 基本类型直接转string
		if httpContext.Route.IsGoBasicType {
			httpContext.Response.BodyString = parse.Convert(responseBody, "")
			httpContext.Response.BodyBytes = []byte(httpContext.Response.BodyString)
		} else { // dto
			httpContext.Response.BodyBytes, _ = json.Marshal(responseBody)
			httpContext.Response.BodyString = string(httpContext.Response.BodyBytes)
		}
		httpContext.Response.StatusCode = 200
		return
	}

	// 多个返回值，则转成数组Json
	lst := collections.NewListAny()
	for i := 0; i < len(httpContext.Response.Body); i++ {
		lst.Add(httpContext.Response.Body[i].Interface())
	}
	httpContext.Response.BodyBytes, _ = json.Marshal(lst)
	httpContext.Response.BodyString = string(httpContext.Response.BodyBytes)
	httpContext.Response.StatusCode = 200
}
