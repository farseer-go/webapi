package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRouteRegexp(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)

	webapi.RegisterGET("/mini/{pageSize}-{pageIndex}", func(pageSize int, pageIndex int) (int, int) {
		return pageSize, pageIndex
	})
	webapi.RegisterPOST("/mini/{pageSize}/{pageIndex}", func(pageSize int, pageIndex int) (int, int) {
		return pageSize, pageIndex
	})

	webapi.UseApiResponse()

	go webapi.Run(":8091")
	time.Sleep(10 * time.Millisecond)

	t.Run("/mini/{pageSize}-{pageIndex}-get", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8091/mini/15-6", nil)
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{15, 6}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("/mini/{pageSize}/{pageIndex}-post", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8091/mini/15/6", nil)
		rsp, _ := http.DefaultClient.Do(req)
		apiResponse := core.NewApiResponseByReader[[]int](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, []int{15, 6}, apiResponse.Data)
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
