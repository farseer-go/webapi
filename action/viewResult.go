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
	action := httpContext.URI.Path[strings.LastIndex(httpContext.URI.Path, "/")+1:]
	path, _ := receiver.cutPrefix(httpContext.URI.Path, "/")
	path = path[:len(path)-len(action)]

	if receiver.viewName == "" {
		receiver.viewName = "./views/" + path + action + ".html"
	} else {
		receiver.viewName = strings.TrimPrefix(receiver.viewName, "/")
		lstViewPath := collections.NewList(strings.Split(receiver.viewName, "/")...)
		if !strings.Contains(lstViewPath.Last(), ".") {
			receiver.viewName = "./views/" + path + receiver.viewName + ".html"
		} else {
			receiver.viewName = "./views/" + path + receiver.viewName
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

	httpContext.Response.Write(file)
}

// CutPrefix go 1.20
func (receiver ViewResult) cutPrefix(s, prefix string) (after string, found bool) {
	if !strings.HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
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
