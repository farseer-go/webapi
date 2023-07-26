package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"testing"
	"time"
)

func TestResponse(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.UseApiResponse()
	go webapi.Run(":8086")
	time.Sleep(10 * time.Millisecond)
}
