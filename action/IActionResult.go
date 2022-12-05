package action

import "github.com/farseer-go/webapi/context"

type IResult interface {
	ExecuteResult(httpContext *context.HttpContext)
}
