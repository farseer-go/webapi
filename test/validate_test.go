package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/webapi"
	"github.com/stretchr/testify/assert"
)

type ValidateRequest struct {
	Name string `validate:"required" label:"账号"`
	Age  int    `validate:"gte=0,lte=100" label:"年龄"`
}

func TestValidate(t *testing.T) {
	fs.Initialize[webapi.Module]("demo")
	configure.SetDefault("Log.Component.webapi", true)
	webapi.RegisterRoutes(webapi.Route{Url: "/Validate", Method: "POST", Action: func(dto ValidateRequest) string {
		return fmt.Sprintf("%+v", dto)
	}})
	webapi.UseValidate()
	go webapi.Run(":8092")
	time.Sleep(10 * time.Millisecond)

	t.Run("/Validate error", func(t *testing.T) {
		sizeRequest := ValidateRequest{Name: "", Age: 200}
		marshal, _ := snc.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8092/Validate", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "账号为必填字段,年龄必须小于或等于100", string(body))
		assert.Equal(t, 403, rsp.StatusCode)
	})

	t.Run("/Validate success", func(t *testing.T) {
		sizeRequest := ValidateRequest{Name: "steden", Age: 37}
		marshal, _ := snc.Marshal(sizeRequest)
		rsp, _ := http.Post("http://127.0.0.1:8092/Validate", "application/json", bytes.NewReader(marshal))
		body, _ := io.ReadAll(rsp.Body)
		_ = rsp.Body.Close()
		assert.Equal(t, "{Name:steden Age:37}", string(body))
		assert.Equal(t, 200, rsp.StatusCode)
	})
}
