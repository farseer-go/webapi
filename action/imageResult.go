package action

import (
	"bytes"
	"github.com/farseer-go/webapi/context"
)

// ImageResult 处理图片
type ImageResult struct {
	buffer    *bytes.Buffer
	imageType string
}

func (receiver *ImageResult) ExecuteResult(httpContext *context.HttpContext) {
	httpContext.Response.Write(receiver.buffer.Bytes())
	httpContext.Response.SetHeader("Content-Type", receiver.imageType)
}

// Image 返回图片格式
func Image(buffer *bytes.Buffer, imageType string) IResult {
	return &ImageResult{
		buffer:    buffer,
		imageType: imageType,
	}
}
