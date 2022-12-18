package action

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
	"os"
	"strings"
)

// ViewResult 返回视图功能
type ViewResult struct {
	ViewName string
}

func (receiver ViewResult) ExecuteResult(httpContext *context.HttpContext) {
	// 默认视图，则以routeUrl为视图位置
	if receiver.ViewName == "" {
		receiver.ViewName = "./views/" + strings.TrimPrefix(httpContext.URI.Path, "/") + ".html"
	} else {
		receiver.ViewName = strings.TrimPrefix(receiver.ViewName, "/")
		lstViewPath := collections.NewList(strings.Split(receiver.ViewName, "/")...)
		if !strings.Contains(lstViewPath.Last(), ".") {
			receiver.ViewName = "./views/" + receiver.ViewName + ".html"
		} else {
			receiver.ViewName = "./views/" + receiver.ViewName
		}
	}

	file, _ := os.ReadFile(receiver.ViewName)
	httpContext.Response.BodyString = string(file)
	httpContext.Response.BodyBytes = []byte(httpContext.Response.BodyString)
	httpContext.Response.StatusCode = 200
}

// View 视图
func View(viewName ...string) IResult {
	var view string
	if len(viewName) > 0 {
		view = viewName[0]
	}
	return ViewResult{
		ViewName: view,
	}
}
