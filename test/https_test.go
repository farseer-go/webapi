package test

import (
	"github.com/farseer-go/webapi"
	"testing"
)

func TestHttps(t *testing.T) {
	server := webapi.NewApplicationBuilder()
	server.UseTLS("", "")
	go server.Run(":80443")
}
