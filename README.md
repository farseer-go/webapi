# webapi 概述
> 包：`"github.com/farseer-go/webapi"`

> 模块：`webapi.Module`

- [文档（国内）](https://farseer-go.gitee.io/)
- [文档（国外）](https://farseer-go.github.io/doc/)
- [开源（Github）](https://github.com/farseer-go/fs)

![](https://img.shields.io/github/stars/farseer-go?style=social)
![](https://img.shields.io/github/license/farseer-go/webapi)
![](https://img.shields.io/github/go-mod/go-version/farseer-go/webapi)
![](https://img.shields.io/github/v/release/farseer-go/webapi)
![](https://img.shields.io/github/languages/code-size/farseer-go/webapi)
![](https://img.shields.io/github/directory-file-count/farseer-go/webapi)
![](https://img.shields.io/github/last-commit/farseer-go/webapi)

> 用于快速构建api服务，带来极简、优雅的开发体验。编写api服务时，不需要使用httpRequest、httpResponse等数据结构。

webapi使用了中间件的管道模型编写，让我们加入非业务逻辑时非常简单。

包含两种风格来提供API服务：
- `minimalApi`：动态API风格（直接绑定到逻辑层）
- `Mvc`：Controller、Action风格

> 使用minimalApi时，甚至不需要UI层来提供API服务。

大部份情况下，除了main需要配置webapi路由外，在你的api handle中就是一个普通的func函数，不需要依赖webapi组件。webapi会根据`func函数`的`出入参`来`隐式绑定数据`。

```go
func main() {
	fs.Initialize[webapi.Module]("FOPS")
	webapi.RegisterPOST("/mini/hello1", Hello1)
	webapi.RegisterPOST("/mini/hello3", Hello3, "pageSize", "pageIndex")
	webapi.Run()
}

// 使用结构（DTO）来接收入参
// 返回string
func Hello1(req pageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

// 使用基础参数来接收入参
// 返回pageSizeRequest结构（会自动转成json)
func Hello3(pageSize int, pageIndex int) pageSizeRequest {
    return pageSizeRequest{
        PageSize:  pageSize,
        PageIndex: pageIndex,
    }
}

// 也可以定义一个结构，用于接收参数
type pageSizeRequest struct {
    PageSize  int
    PageIndex int
}
```

函数中，`出入参都会自动绑定数据`

> 如果是`application/json`，则会自动被反序列化成model，如果是`x-www-form-urlencoded`，则会将每一项的key/value匹配到model字段中

可以看到，整个过程，`不需要`做`json序列化`、`httpRequest`、`httpResponse`的操作。