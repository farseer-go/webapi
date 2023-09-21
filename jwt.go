package webapi

import (
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte                  // jwt key
var jwtKeyType string              // key type
var jwtKeyMethod jwt.SigningMethod // sign method

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

// NewJwtToken 生成Jwt Token
func NewJwtToken(claims map[string]any) (string, error) {
	//jwt.MapClaims{
	//	"foo": "bar",
	//	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	//}
	token := jwt.NewWithClaims(jwtKeyMethod, jwt.MapClaims(claims))
	if len(jwtKey) == 0 {
		return token.SigningString()
	}
	return token.SignedString(jwtKey)
}
