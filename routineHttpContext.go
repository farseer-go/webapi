package webapi

import (
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/webapi/context"
)

var routineHttpContext = asyncLocal.New[*context.HttpContext]()
