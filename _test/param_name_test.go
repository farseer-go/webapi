package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"testing"
)

func hello(pageSize int, pageIndex int) (int, int) {
	return pageSize, pageIndex
}

func TestParamName(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.RegisterDELETE("/cors/test", hello)
	webapi.UseCors()
	go webapi.Run(":8080")
}
