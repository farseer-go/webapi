package context

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// HeaderName 前端提交Token，存放到header的name
var HeaderName string
var JwtKey []byte                  // 生成token的秘钥
var JwtKeyType string              // 签名加密方式
var InvalidStatusCode int          // token无效时的状态码
var InvalidMessage string          // key type
var JwtKeyMethod jwt.SigningMethod // sign method

type HttpJwt struct {
	r      *http.Request
	w      http.ResponseWriter
	claims jwt.MapClaims
}

// GetToken 获取前端提交过来的Token
func (receiver *HttpJwt) GetToken() string {
	return receiver.r.Header.Get(HeaderName)
}

// Build 生成jwt token，并写入head
func (receiver *HttpJwt) Build() (string, error) {
	claims := make(map[string]any)
	claims["farseer-go"] = "v0.8.0"

	// 生成token对象
	token := jwt.NewWithClaims(JwtKeyMethod, jwt.MapClaims(claims))
	var sign string
	var err error

	if len(JwtKey) == 0 {
		sign, err = token.SigningString() // 不带秘钥的签名
	} else {
		sign, err = token.SignedString(JwtKey) // 带秘钥的签名
	}

	// 成功生成后，写入到head
	if err == nil {
		receiver.w.Header().Set(HeaderName, sign)
	}
	return sign, err
}

// Valid 验证前端提交过来的token是否正确
func (receiver *HttpJwt) Valid() bool {
	token, err := jwt.Parse(receiver.GetToken(), func(token *jwt.Token) (any, error) {
		// 验证加密方式
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法： %v", token.Header["alg"])
		}
		return JwtKey, nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		receiver.claims = claims
		return true
	}
	return false
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}
