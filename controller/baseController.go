package controller

import (
	"github.com/farseer-go/webapi/context"
	"net/http"
)

type BaseController struct {
	HttpContext context.HttpContext
}

func (receiver BaseController) init(r *http.Request) {

}
