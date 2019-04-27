package JWT

import (
	"testing"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsIm5hbWUiOiLmtYvor5XnlKjmiLcxIiwicGhvbmUiOiJ0aXB1X3RlbCIsImVtYWlsIjoidGVzdEB0aXB1LmNvbSIsImV4cCI6MTU1NjM0NzcxOH0.R4FxKyVLbJK6T_u1Qk0J1bbhvHvdNkGtUhO6RWBFDao"

var j = NewJWT()

// 测试生成令牌
func TestJWT_CreateToken(t *testing.T) {
	token, err := j.CreateToken(1, "测试用户1")

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
	claims, err := j.RefreshToken(token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(claims)
}
