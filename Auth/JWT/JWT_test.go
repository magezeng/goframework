package JWT

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwibmFtZSI6Iua1i-ivleeUqOaItzEiLCJwaG9uZSI6IjEzODQ4MzgyMzE4IiwiZW1haWwiOiI1MDMzMDk5NEBxcS5jb20iLCJleHAiOjE1NTIyMjU1MTksImlzcyI6IiMjIyMjVGlwdSFAIzEyMyIsIm5iZiI6MTU1MjIyMDkxOX0.-LLJ6TRxflsDx6BeFZr5IIvmUFP-i-_MQ0_vcy_dgC0"

var j = NewJWT()

// 测试生成令牌
func TestJWT_CreateToken(t *testing.T) {
	claims := CustomClaims{
		"1",
		"测试用户1",
		"13848382318",
		"50330994@qq.com",
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("生成完成，token为: ", token)
}

// 测试解析令牌
func TestJWT_ParseToken(t *testing.T) {
	claims, err := j.ParseToken(token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(claims)
}

// 测试刷新令牌
func TestJWT_RefreshToken(t *testing.T) {
	claims, err := j.ParseToken(token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(claims)
}
