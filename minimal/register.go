package minimal

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/types"
	"github.com/farseer-go/webapi/check"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/middleware"
	"reflect"
	"strings"
)

// Register 注册单个Api
func Register(area string, method string, route string, actionFunc any, filters []context.IFilter, paramNames ...string) *context.HttpRoute {
	actionType := reflect.TypeOf(actionFunc)
	params := types.GetInParam(actionType)

	// 如果设置了方法的入参（多参数），则需要全部设置
	if len(paramNames) > 0 && len(paramNames) != len(params) {
		flog.Panicf("注册路由%s%s失败：%s函数入参与%s不匹配，建议重新运行fsctl -r命令", area, route, flog.Red(actionType.String()), flog.Blue(paramNames))
	}

	lstRequestParamType := collections.NewList(params...)
	lstResponseParamType := collections.NewList(types.GetOutParam(actionType)...)

	// 入参是否为DTO模式
	isDtoModel := types.IsDtoModelIgnoreInterface(params)
	// 是否实现了ICheck
	var requestParamIsImplCheck bool
	if isDtoModel {
		// 是否实现了check.ICheck
		var checker = reflect.TypeOf((*check.ICheck)(nil)).Elem()
		requestParamIsImplCheck = lstRequestParamType.First().Implements(checker)
		if !requestParamIsImplCheck {
			requestParamIsImplCheck = reflect.PointerTo(lstRequestParamType.First()).Implements(checker)
		}
	}

	// 添加到路由表
	return &context.HttpRoute{
		Schema:                  "http",
		RouteUrl:                area + route,
		Action:                  actionFunc,
		Method:                  collections.NewList(strings.Split(strings.ToUpper(method), "|")...),
		RequestParamType:        lstRequestParamType,
		ResponseBodyType:        lstResponseParamType,
		RequestParamIsImplCheck: requestParamIsImplCheck,
		RequestParamIsModel:     isDtoModel,
		ResponseBodyIsModel:     types.IsDtoModel(lstResponseParamType.ToArray()),
		ParamNames:              collections.NewList(paramNames...),
		HttpMiddleware:          &middleware.Http{},
		HandleMiddleware:        &HandleMiddleware{},
		IsGoBasicType:           types.IsGoBasicType(lstResponseParamType.First()),
		Filters:                 filters,
	}
}
