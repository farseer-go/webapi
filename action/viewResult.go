package action

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/context"
	"os"
	"strings"
)

// ViewResult 返回视图功能
type ViewResult struct {
	viewName string
	data     map[string]any
}

func (receiver ViewResult) ExecuteResult(httpContext *context.HttpContext) {
	// 默认视图，则以routeUrl为视图位置
	if receiver.viewName == "" {
		receiver.viewName = "./views/" + strings.TrimPrefix(httpContext.URI.Path, "/") + ".html"
	} else {
		receiver.viewName = strings.TrimPrefix(receiver.viewName, "/")
		lstViewPath := collections.NewList(strings.Split(receiver.viewName, "/")...)
		if !strings.Contains(lstViewPath.Last(), ".") {
			receiver.viewName = "./views/" + receiver.viewName + ".html"
		} else {
			receiver.viewName = "./views/" + receiver.viewName
		}
	}

	file, _ := os.ReadFile(receiver.viewName)
	//if len(receiver.data) > 0 {
	//	htmlSource := string(file)
	//	startIndex := -1
	//	stopIndex := -1
	//	// 遍历html字符串
	//	for i := 0; i < len(htmlSource); i++ {
	//		// 查找开始标记
	//		if i < len(htmlSource)-1 {
	//			if htmlSource[i:1] == "${" {
	//				startIndex = i
	//				i++
	//			}
	//		}
	//
	//		// 查找结束标记
	//		if startIndex > -1 && htmlSource[i] == '}' {
	//			stopIndex = i
	//
	//			// 找到了
	//			htmlSource[startIndex : stopIndex-startIndex]
	//		}
	//	}
	//	for k, v := range receiver.data {
	//		htmlSource = strings.ReplaceAll(htmlSource, "${"+k+"}", parse.Convert(v, ""))
	//	}
	//} else {
	//
	//}
	httpContext.Response.BodyString = string(file)
	httpContext.Response.BodyBytes = file
	httpContext.Response.StatusCode = 200
}

// View 视图
func View(viewName ...string) IResult {
	var view string
	if len(viewName) > 0 {
		view = viewName[0]
	}
	return ViewResult{
		viewName: view,
	}
}

// ViewData 视图
func ViewData(Data map[string]any, viewName ...string) IResult {
	var view string
	if len(viewName) > 0 {
		view = viewName[0]
	}
	return ViewResult{
		viewName: view,
		data:     Data,
	}
}
