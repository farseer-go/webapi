package test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
)

func TestArea(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.Area("api/1.0", func() {
		webapi.RegisterGET("/mini/test", func() string {
			return "ok"
		})
	})
	go webapi.Run(":8087")
	time.Sleep(10 * time.Millisecond)

	rsp, _ := http.Get("http://127.0.0.1:8087/api/1.0/mini/test")
	body, _ := io.ReadAll(rsp.Body)
	_ = rsp.Body.Close()
	assert.Equal(t, "ok", string(body))
}
