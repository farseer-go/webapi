package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func init() {
	fs.Initialize[webapi.Module]("demo")
}
