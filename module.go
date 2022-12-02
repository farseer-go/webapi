package webapi

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/minimal"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	context.LstRouteTable = collections.NewList[context.HttpRoute]()

	controller.Init()
	minimal.Init()

	middleware.Init()
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
