package test

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGetHttpContext(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")

	webapi.RegisterPOST("/getHttpContext", func() string {
		return webapi.GetHttpContext().ContentType
	})
	webapi.UseApiResponse()
	go webapi.Run(":8089")
	time.Sleep(10 * time.Millisecond)

	t.Run("getHttpContext", func(t *testing.T) {
		rsp, _ := http.Post("http://127.0.0.1:8089/getHttpContext", "application/json", nil)
		apiResponse := core.NewApiResponseByReader[string](rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "application/json", apiResponse.Data)
	})
}
