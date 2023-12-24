package session_redis

import (
	webapiContext "github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/session-redis/context"
	"github.com/farseer-go/webapi/session-redis/middleware"
)

type SessionMiddlewareCreat struct {
}

// Create 创建Session中间件
func (receiver *SessionMiddlewareCreat) Create() webapiContext.IMiddleware {
	go context.ClearSession()
	return &middleware.Session{}
}
