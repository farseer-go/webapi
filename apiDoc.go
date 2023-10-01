package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/action"
	"github.com/farseer-go/webapi/context"
	"reflect"
	"strings"
)

// UseApiDoc 是否开启Api文档
func (r *applicationBuilder) UseApiDoc() {
	r.registerAction(Route{Url: "/doc/api", Method: "GET", Action: func() action.IResult {
		lstBody := collections.NewList[string]("<html><body>\n")
		lstRoute := collections.NewList[context.HttpRoute]()
		for _, httpRoute := range r.mux.m {
			lstRoute.Add(*httpRoute)
		}

		// 遍历路由
		lstRoute.OrderBy(func(item context.HttpRoute) any {
			return item.RouteUrl
		}).Foreach(func(httpRoute *context.HttpRoute) {
			method := strings.Join(httpRoute.Method.ToArray(), "|")
			if httpRoute.RouteUrl == "/" && method == "" && httpRoute.Controller == nil && httpRoute.Action == nil {
				return
			}
			// API地址
			lstBody.Add(fmt.Sprintf("[%s]：<a href=\"%s\" target=\"_blank\">%s</a><br />\n", method, r.hostAddress+httpRoute.RouteUrl, r.hostAddress+httpRoute.RouteUrl))
			// 入参
			// 使用textarea
			lstBody.Add("<textarea style=\"width: 50%; height: 120px;\">")
			// dto模式转成json
			if httpRoute.RequestParamIsModel {
				val := reflect.New(httpRoute.RequestParamType.First()).Interface()
				indent, _ := json.MarshalIndent(val, "", "  ")
				lstBody.Add(string(indent))
			} else {
				var mapVal = make(map[string]any)
				for i := 0; i < httpRoute.RequestParamType.Count(); i++ {
					fieldType := httpRoute.RequestParamType.Index(i)
					// 必须是非interface类型
					if fieldType.Kind() != reflect.Interface && httpRoute.ParamNames.Count() > i {
						// 指定了参数名称
						paramName := strings.ToLower(httpRoute.ParamNames.Index(i))
						mapVal[paramName] = reflect.New(fieldType).Elem().Interface()
					}
				}
				indent, _ := json.MarshalIndent(mapVal, "", "  ")
				lstBody.Add(string(indent))
			}
			lstBody.Add("</textarea>")

			lstBody.Add("<textarea style=\"width: 50%; height: 120px;\">")
			// 只有一个返回值
			if httpRoute.ResponseBodyType.Count() == 1 {
				responseBody := reflect.New(httpRoute.ResponseBodyType.First()).Interface()
				// 基本类型直接转string
				if httpRoute.IsGoBasicType {
					lstBody.Add(parse.ToString(responseBody))
				} else { // dto
					indent, _ := json.MarshalIndent(responseBody, "", "  ")
					lstBody.Add(string(indent))
				}
			} else {
				// 多个返回值，则转成数组Json
				lst := collections.NewListAny()
				httpRoute.ResponseBodyType.Foreach(func(item *reflect.Type) {
					lst.Add(reflect.New(*item).Interface())
				})
				indent, _ := json.MarshalIndent(lst, "", "  ")
				lstBody.Add(string(indent))
			}
			lstBody.Add("</textarea>")
			lstBody.Add("<hr />")
		})
		lstBody.Add("</body><html>")

		return action.Content(lstBody.ToString(""))
	}, Params: nil})
}
