package webapi

import (
	"github.com/farseer-go/fs"
)

func ExampleRun() {
	fs.Initialize[Module]("demo")
	Run()
}
