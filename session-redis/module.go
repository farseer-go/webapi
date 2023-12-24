package session_redis

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/session-redis/context"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	sessionTimeout := configure.GetInt("Webapi.Session.Age")
	if sessionTimeout > 0 {
		context.SessionTimeout = sessionTimeout
	}
}
