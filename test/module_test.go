package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
	"testing"
)

func TestModule(t *testing.T) {
	fs.Initialize[webapi.Module]("unit test")
	webapi.Module{}.Shutdown()
}
