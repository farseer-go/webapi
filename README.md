# webapi 启动
> [文档：https://farseer-go.github.io/doc/](https://farseer-go.github.io/doc/)

> 包：`"github.com/farseer-go/webapi"`
>
> 模块：`webapi.Module`

## 启动Web服务
_startupModule.go 启动模块_
```go
package main
import (
    "github.com/farseer-go/fs/modules"
    "github.com/farseer-go/webapi"
)
type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}}
}

...
```
!> 这里依赖了`webapi.Module`模块，这是必须的。

_main.go 入口_
```go
package main
import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("FOPS")
	webapi.RegisterPOST("/mini/hello1", Hello1)
	webapi.Run()
}

func Hello1(req request.PageSizeRequest) string {
	return fmt.Sprintf("hello world pageSize=%d，pageIndex=%d", req.PageSize, req.PageIndex)
}

type PageSizeRequest struct {
	PageSize  int
	PageIndex int
}
```

执行`webapi.Run()`后，便可启动web api服务。

_运行结果：_

```
2022-12-01 17:07:24 基础组件初始化完成
2022-12-01 17:07:24 初始化完毕，共耗时：1 ms 
2022-12-01 17:07:24 ---------------------------------------
2022-12-01 17:07:24 [Info] Web服务已启动：http://localhost:8888/
```

> 在Hello1函数中，`不需要依赖任何框架的参数`，这是`webapi的特色`之一：`极简`、`优雅`
>
> 我们甚至不需要`接口层（UI层）`，可以直接通过路由指向`逻辑层（或应用层）`

## 端口

在不做任何配置时，EndPort默认为：`localhost:8888`

配置EndPort的两种方式：

### 1、参数配置
```go
func Run(params ...string)
```
params实际只支持一个入参，支持的格式为:
- webapi.Run(`""`)
- webapi.Run(`":80"`)
- webapi.Run(`"127.0.0.1:80"`)

### 2、配置文件
_./farseer.yaml_
```yaml
WebApi:
  Url: ":8888"
```
?> 支持的格式与参数配置是一样的。