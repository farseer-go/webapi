package test

import (
	"bytes"
	"encoding/json"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNewApplicationBuilder(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	server := webapi.NewApplicationBuilder()
	server.RegisterController(&TestController{})
	server.RegisterPOST("/mini/hello1", Hello1)
	server.RegisterDELETE("/mini/hello4", Hello4, "pageSize", "pageIndex")
	server.RegisterPOST("/mini/hello5", Hello5)
	server.RegisterPOST("/mini/hello6", Hello6)
	server.RegisterPOST("/mini/hello8", Hello8)
	server.UseCors()
	go server.Run(":8889")
	time.Sleep(100 * time.Millisecond)

	t.Run("mini/hello5:8889", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello5", "application/json", nil)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "s501", string(body))
		assert.Equal(t, 501, rsp.StatusCode)
	})

	t.Run("mini/hello6:8889", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello6", "application/json", nil)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "s500", string(body))
		assert.Equal(t, 500, rsp.StatusCode)
	})
	t.Run("mini/hello1:8889", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello1", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "hello world pageSize=10，pageIndex=2", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello8:8889", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8889/mini/hello8", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello4:8889", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8889/mini/hello4", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "[10,2]", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("mini/hello4:8889-OPTIONS", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("OPTIONS", "http://127.0.0.1:8889/mini/hello4", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "localhost")
		rsp, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "", string(body))
		assert.Equal(t, "localhost", rsp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", rsp.Header.Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, 204, rsp.StatusCode)
	})
}
