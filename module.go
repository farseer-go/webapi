package webapi

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/controller"
	"github.com/farseer-go/webapi/minimal"
	"github.com/golang-jwt/jwt/v5"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	controller.Init()
	minimal.Init()
	defaultApi = NewApplicationBuilder()

	sessionTimeout := configure.GetInt("Webapi.Session.Age")
	if sessionTimeout > 0 {
		context.SessionTimeout = sessionTimeout
	}

	// jwt
	context.JwtKey = []byte(configure.GetString("WebApi.Jwt.Key"))
	context.JwtKeyType = configure.GetString("WebApi.Jwt.KeyType")
	context.HeaderName = configure.GetString("WebApi.Jwt.HeaderName")
	context.InvalidStatusCode = configure.GetInt("WebApi.Jwt.InvalidStatusCode")
	context.InvalidMessage = configure.GetString("WebApi.Jwt.InvalidMessage")

	if context.HeaderName == "" {
		context.HeaderName = "Authorization"
	}

	if context.InvalidStatusCode == 0 {
		context.InvalidStatusCode = 403
	}

	if context.InvalidMessage == "" {
		context.InvalidMessage = "您没有权限访问"
	}

	switch context.JwtKeyType {
	case "HS256":
		context.JwtKeyMethod = jwt.SigningMethodHS256
	case "HS384":
		context.JwtKeyMethod = jwt.SigningMethodHS384
	case "HS512":
		context.JwtKeyMethod = jwt.SigningMethodHS512

	case "RS256":
		context.JwtKeyMethod = jwt.SigningMethodRS256
	case "RS384":
		context.JwtKeyMethod = jwt.SigningMethodRS384
	case "RS512":
		context.JwtKeyMethod = jwt.SigningMethodRS512

	case "ES256":
		context.JwtKeyMethod = jwt.SigningMethodES256
	case "ES384":
		context.JwtKeyMethod = jwt.SigningMethodES384
	case "ES512":
		context.JwtKeyMethod = jwt.SigningMethodES512

	case "PS256":
		context.JwtKeyMethod = jwt.SigningMethodPS256
	case "PS384":
		context.JwtKeyMethod = jwt.SigningMethodPS384
	case "PS512":
		context.JwtKeyMethod = jwt.SigningMethodPS512

	case "EdDSA":
		context.JwtKeyMethod = jwt.SigningMethodEdDSA
	default:
		context.JwtKeyMethod = jwt.SigningMethodHS256
	}
}
