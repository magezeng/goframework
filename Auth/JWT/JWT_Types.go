package JWT

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID    uint `json:"userId"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("令牌已经过期")
	TokenNotValidYet = errors.New("令牌未被激活")
	TokenMalformed   = errors.New("该令牌不是由tipu发出的")
	TokenInvalid     = errors.New("令牌无法被解析")
	SecretKey          = "Tipu!@#123"
)
