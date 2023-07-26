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

func TestCors(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	server := webapi.NewApplicationBuilder()
	server.RegisterDELETE("/cors/test", func(pageSize int, pageIndex int) (int, int) {
		return pageSize, pageIndex
	}, "page_Size", "pageIndex")
	server.UseCors()
	go server.Run(":8080")
	time.Sleep(100 * time.Millisecond)

	t.Run("/cors/test:8080", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8080/cors/test", bytes.NewReader(marshal))
		req.Header.Set("Content-Type", "application/json")
		rsp, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "[10,2]", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("/cors/test:8080-OPTIONS", func(t *testing.T) {
		sizeRequest := pageSizeRequest{PageSize: 10, PageIndex: 2}
		marshal, _ := json.Marshal(sizeRequest)
		req, _ := http.NewRequest("OPTIONS", "http://127.0.0.1:8080/cors/test", bytes.NewReader(marshal))
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
