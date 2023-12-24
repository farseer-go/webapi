package webapi

import (
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/minimal"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	controller.Init()
	minimal.Init()
	defaultApi = NewApplicationBuilder()
}

func (module Module) PostInitialize() {
	context.InitJwt()
}
