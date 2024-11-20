package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
)

func TestNewApplicationBuilder(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	server := webapi.NewApplicationBuilder()
	server.RegisterPOST("/mini/test", func() {})
	go server.Run(":8083")
	time.Sleep(10 * time.Millisecond)

	t.Run("mini/test2:8083", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := sonic.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8083/mini/test", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
