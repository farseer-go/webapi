package webapi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/webapi/action"
	"github.com/farseer-go/webapi/context"
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
				indent, _ := sonic.MarshalIndent(val, "", "  ")
				lstBody.Add(string(indent))
			} else {
				var mapVal = make(map[string]any)
				for i := 0; i < httpRoute.RequestParamType.Count(); i++ {
					fieldType := httpRoute.RequestParamType.Index(i)
					// 必须是非interface类型
					if fieldType.Kind() != reflect.Interface && httpRoute.ParamNames.Count() > i {
						// 指定了参数名称
						paramName := httpRoute.ParamNames.Index(i)
						mapVal[paramName] = reflect.New(fieldType).Elem().Interface()
					}
				}
				indent, _ := sonic.MarshalIndent(mapVal, "", "  ")
				lstBody.Add(string(indent))
			}
			lstBody.Add("</textarea>")

			lstBody.Add("<textarea style=\"width: 50%; height: 120px;\">")
			if httpRoute.ResponseBodyType.Count() == 0 { // 没有返回值
				// 使用ApiResponse，则返回ApiResponse格式
				if r.useApiResponse {
					rsp := core.Success[any]("成功", nil)
					indent, _ := sonic.MarshalIndent(rsp, "", "  ")
					lstBody.Add(string(indent))
				} else {
					lstBody.Add("")
				}
			} else if httpRoute.ResponseBodyType.Count() == 1 { // 只有一个返回值
				rspRefType := httpRoute.ResponseBodyType.First()
				responseBody := reflect.New(rspRefType).Interface()
				// 如果是list，则添加3个元素，用于演示
				pointerMeta := fastReflect.PointerOf(responseBody)
				isList := pointerMeta.Type == fastReflect.List
				if isList {
					val := pointerMeta.GetItemMeta().ZeroValue
					lstVal := types.ListNew(pointerMeta.ReflectType)
					types.ListAdd(lstVal, val)
					types.ListAdd(lstVal, val)
					types.ListAdd(lstVal, val)
					responseBody = lstVal.Interface()
				}
				// 使用ApiResponse，则返回ApiResponse格式
				if r.useApiResponse {
					rsp := core.Success[any]("成功", responseBody)
					indent, _ := sonic.MarshalIndent(rsp, "", "  ")
					lstBody.Add(string(indent))
				} else {
					// 基本类型直接转string
					if httpRoute.IsGoBasicType {
						lstBody.Add(parse.ToString(responseBody))
					} else { // dto
						indent, _ := sonic.MarshalIndent(responseBody, "", "  ")
						lstBody.Add(string(indent))
					}
				}
			} else {
				// 多个返回值，则转成数组Json
				lst := collections.NewListAny()
				httpRoute.ResponseBodyType.Foreach(func(item *reflect.Type) {
					lst.Add(reflect.New(*item).Interface())
				})
				indent, _ := sonic.MarshalIndent(lst, "", "  ")
				lstBody.Add(string(indent))
			}
			lstBody.Add("</textarea>")
			lstBody.Add("<hr />")
		})
		lstBody.Add("</body><html>")

		return action.Content(lstBody.ToString(""))
	}, Params: nil})
}
