package action

import (
	"github.com/farseer-go/webapi/context"
	"os"
)

// FileContentResult 返回文件内容
type FileContentResult struct {
	filePath string
}

func (receiver FileContentResult) ExecuteResult(httpContext *context.HttpContext) {
	file, _ := os.ReadFile(receiver.filePath)
	httpContext.Response.Write(file)
	httpContext.Response.StatusCode = 200
}

// FileContent 文件
func FileContent(filePath string) IResult {
	return FileContentResult{
		filePath: filePath,
	}
}
