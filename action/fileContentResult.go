package action

import (
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/webapi/context"
)

// FileContentResult 返回文件内容
type FileContentResult struct {
	filePath string
}

func (receiver FileContentResult) ExecuteResult(httpContext *context.HttpContext) {
	httpContext.Response.BodyString = file.ReadString(receiver.filePath)
	httpContext.Response.BodyBytes = []byte(httpContext.Response.BodyString)
	httpContext.Response.StatusCode = 200
}

// FileContent 文件
func FileContent(filePath string) IResult {
	return FileContentResult{
		filePath: filePath,
	}
}
