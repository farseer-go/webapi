package context

import (
	"fmt"
	"net/http"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/exception"
	"github.com/golang-jwt/jwt/v5"
)

var headerName string              // 前端提交Token，存放到header的name
var jwtKey []byte                  // 生成token的秘钥
var jwtKeyType string              // 签名加密方式
var jwtKeyMethod jwt.SigningMethod // sign method
var InvalidStatusCode int          // token无效时的状态码
var InvalidMessage string          // key type

type HttpJwt struct {
	r      *http.Request
	w      http.ResponseWriter
	claims jwt.MapClaims
}

// InitJwt jwt初始化
func InitJwt() {
	jwtKey = []byte(configure.GetString("WebApi.Jwt.Key"))
	jwtKeyType = configure.GetString("WebApi.Jwt.KeyType")
	headerName = configure.GetString("WebApi.Jwt.HeaderName")
	InvalidStatusCode = configure.GetInt("WebApi.Jwt.InvalidStatusCode")
	InvalidMessage = configure.GetString("WebApi.Jwt.InvalidMessage")

	if headerName == "" {
		headerName = "Authorization"
	}

	if InvalidStatusCode == 0 {
		InvalidStatusCode = 403
	}

	if InvalidMessage == "" {
		InvalidMessage = "您没有权限访问"
	}

	switch jwtKeyType {
	case "HS256":
		jwtKeyMethod = jwt.SigningMethodHS256
	case "HS384":
		jwtKeyMethod = jwt.SigningMethodHS384
	case "HS512":
		jwtKeyMethod = jwt.SigningMethodHS512

	case "RS256":
		jwtKeyMethod = jwt.SigningMethodRS256
	case "RS384":
		jwtKeyMethod = jwt.SigningMethodRS384
	case "RS512":
		jwtKeyMethod = jwt.SigningMethodRS512

	case "ES256":
		jwtKeyMethod = jwt.SigningMethodES256
	case "ES384":
		jwtKeyMethod = jwt.SigningMethodES384
	case "ES512":
		jwtKeyMethod = jwt.SigningMethodES512

	case "PS256":
		jwtKeyMethod = jwt.SigningMethodPS256
	case "PS384":
		jwtKeyMethod = jwt.SigningMethodPS384
	case "PS512":
		jwtKeyMethod = jwt.SigningMethodPS512

	case "EdDSA":
		jwtKeyMethod = jwt.SigningMethodEdDSA
	default:
		jwtKeyMethod = jwt.SigningMethodHS256
	}
}

// GetToken 获取前端提交过来的Token
func (receiver *HttpJwt) GetToken() string {
	// 走ws协议
	if receiver.r.Header.Get("Upgrade") == "websocket" {
		return receiver.r.Form.Get(headerName)
	}
	return receiver.r.Header.Get(headerName)
}

// Build 生成jwt token，并写入head
func (receiver *HttpJwt) Build(claims map[string]any) (string, error) {
	// 生成token对象
	token := jwt.NewWithClaims(jwtKeyMethod, jwt.MapClaims(claims))

	exception.ThrowWebExceptionBool(len(jwtKey) == 0, 403, "未设置jwt的秘钥，请到./farseer.yaml的WebApi.Jwt.Key节点中配置")

	sign, err := token.SignedString(jwtKey) // 带秘钥的签名

	// 成功生成后，写入到head
	if err == nil && receiver.w != nil {
		receiver.w.Header().Set(headerName, sign)
	}
	return sign, err
}

// Valid 验证前端提交过来的token是否正确
func (receiver *HttpJwt) Valid() bool {
	t := receiver.GetToken()
	token, err := jwt.Parse(t, func(token *jwt.Token) (any, error) {
		// 验证加密方式
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法： %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	// 签名不对
	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		receiver.claims = claims
		return true
	}
	return false
}

// GetClaims 读取前端提交过来的Claims
func (receiver *HttpJwt) GetClaims() jwt.MapClaims {
	// nil说明没有执行过验证
	if receiver.claims == nil {
		receiver.Valid()
	}
	return receiver.claims
}

// 清除token
func (receiver *HttpJwt) Clear() {
	receiver.w.Header().Del(headerName)
}

//type MyCustomClaims struct {
//	Foo string `json:"foo"`
//	jwt.RegisteredClaims
//}
