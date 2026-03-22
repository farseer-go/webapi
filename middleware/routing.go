package middleware

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/farseer-go/webapi/context"
	"github.com/klauspost/compress/zstd"
)

// zstdDecoder 包级解码器，并发安全，复用以避免重复初始化开销
var zstdDecoder, _ = zstd.NewReader(nil)

type routing struct {
	context.IMiddleware
}

func (receiver *routing) Invoke(httpContext *context.HttpContext) {
	// 检查method
	if httpContext.Route.Schema != "ws" && httpContext.Method != "OPTIONS" && !httpContext.Route.Method.Contains(httpContext.Method) {
		// 响应码
		httpContext.Response.Reject(405, "405 Method NotAllowed")
		return
	}

	// 解析Body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(httpContext.Request.R.Body)
	// 智能解压：根据 Content-Encoding 头自动选择解压算法
	httpContext.Request.BodyBytes = decompressBody(buf.Bytes(), httpContext.Request.R.Header.Get("Content-Encoding"))

	// 解析请求的参数
	httpContext.Request.ParseQuery()
	httpContext.Request.ParseForm()
	httpContext.URI.Query = httpContext.Request.Query

	// 转换成Handle函数需要的参数
	httpContext.Request.Params = httpContext.ParseParams()

	receiver.IMiddleware.Invoke(httpContext)
}

// decompressBody 根据 Content-Encoding 对 body 进行解压。
// 支持：zstd、gzip、deflate、br（brotli）。
// 解压失败或编码未知时原样返回原始数据。
func decompressBody(data []byte, encoding string) []byte {
	if len(data) == 0 || encoding == "" {
		return data
	}
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "zstd":
		result, err := zstdDecoder.DecodeAll(data, make([]byte, 0, len(data)*3))
		if err != nil {
			return data
		}
		return result

	case "gzip":
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return data
		}
		defer r.Close()
		result, err := io.ReadAll(r)
		if err != nil {
			return data
		}
		return result

	case "deflate":
		r := flate.NewReader(bytes.NewReader(data))
		defer r.Close()
		result, err := io.ReadAll(r)
		if err != nil {
			return data
		}
		return result

	case "br":
		r := brotli.NewReader(bytes.NewReader(data))
		result, err := io.ReadAll(r)
		if err != nil {
			return data
		}
		return result
	}
	return data
}
