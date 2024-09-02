package websocket

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/middleware"
	"reflect"
	"strings"
)

// Register 注册单个Api
func Register(area string, method string, route string, actionFunc any, filters []context.IFilter) *context.HttpRoute {
	actionType := reflect.TypeOf(actionFunc)
	inParams := types.GetInParam(actionType)
	outParams := types.GetOutParam(actionType)

	if len(inParams) != 1 || !strings.HasPrefix(inParams[0].String(), "*websocket.Context[") {
		flog.Panicf("注册ws路由%s%s失败：%s函数入参必须为：%s", area, route, flog.Red(actionType.String()), flog.Blue("*websocket.Context[T any]"))
	}
	if len(outParams) != 0 {
		flog.Panicf("注册ws路由%s%s失败：%s函数不能有出参", area, route, flog.Red(actionType.String()))
	}

	// 入参的泛型是否为DTO模式
	itemTypeMethod, _ := inParams[0].MethodByName("ItemType")
	itemType := itemTypeMethod.Type.Out(0)
	isDtoModel := types.IsDtoModelIgnoreInterface([]reflect.Type{itemType})

	// 添加到路由表
	return &context.HttpRoute{
		Schema:              "ws",
		RouteUrl:            area + route,
		Action:              actionFunc,
		Method:              collections.NewList(strings.Split(strings.ToUpper(method), "|")...),
		RequestParamType:    collections.NewList(inParams[0]),
		ResponseBodyType:    collections.NewList[reflect.Type](),
		RequestParamIsModel: isDtoModel,
		ResponseBodyIsModel: false,
		ParamNames:          collections.NewList[string]("context"),
		HttpMiddleware:      &middleware.Websocket{},
		HandleMiddleware:    &HandleMiddleware{},
		IsGoBasicType:       types.IsGoBasicType(itemType),
		Filters:             filters,
	}
}
