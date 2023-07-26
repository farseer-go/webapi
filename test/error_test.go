package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	webapi.RegisterPOST("/error/1", func() {
		exception.ThrowWebException(501, "s501")
	})

	webapi.RegisterPOST("/error/2", func() {
		exception.ThrowException("s500")
	})
	webapi.UseApiResponse()
	go webapi.Run(":8081")
	time.Sleep(100 * time.Millisecond)

	t.Run("error/1", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8081/error/1", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[any](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 501, apiResponse.StatusCode)
		assert.Equal(t, "s501", apiResponse.StatusMessage)
		assert.Equal(t, false, apiResponse.Status)
		assert.Equal(t, 200, rsp.StatusCode)
	})

	t.Run("error/2", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8081/error/2", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[any](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, 500, apiResponse.StatusCode)
		assert.Equal(t, "s500", apiResponse.StatusMessage)
		assert.Equal(t, false, apiResponse.Status)
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
