package action

import (
	"github.com/farseer-go/webapi/context"
)

// ContentResult 返回响应内容
type ContentResult struct {
	content string
}

func (receiver ContentResult) ExecuteResult(httpContext *context.HttpContext) {
	httpContext.Response.WriteString(receiver.content)
}

// Content 内容
func Content(content string) IResult {
	return ContentResult{
		content: content,
	}
}
