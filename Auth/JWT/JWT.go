package JWT

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWT 签名结构
type JWT struct {
	SecretKey []byte
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(SecretKey),
	}
}

// CreateToken 生成一个Token
func (j *JWT) CreateToken(id uint64, name string) (string, error) {
	// 暂时还没有邮箱和电话
	// 邮箱由于发送邮件
	claims := j.createCustomClaims(id, name, "test@tipu.com", "tipu_tel")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SecretKey)
}

func (j *JWT) createCustomClaims(id uint64, name string, email string, phone string) CustomClaims {
	return CustomClaims{
		id,
		name,
		phone,
		email,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
		},
	}
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// RefreshToken 更新Token，增加一个小时的过期时间
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(claims.ID, claims.Name)
	}
	return "", TokenInvalid
}
