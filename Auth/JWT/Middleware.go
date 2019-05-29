package JWT

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"tipu.com/go-framework/Code"
	"tipu.com/go-framework/Models"
)

// Middleware gin中间件，检查token
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// 可能是从websocket过来的，此时就不在Authorization里面
			token = c.Request.Header.Get("Sec-WebSocket-Protocol")
		}
		if token == "" {
			Models.ResultFail(c, Code.UNAUTHORIZED_ERROR, errors.New("无权限访问!"))
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken 解析token包含的信息，判断是否过期，有效性等
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				Models.ResultFail(c, Code.TOKEN_EXPIRED_ERROR, errors.New("Token已经过期!"))
				c.Abort()
				return
			}
			Models.ResultFail(c, Code.TOKEN_ERROR, err)
			c.Abort()
			return
		}
		c.Set("token", token)
		c.Set("claims", claims)
	}
}
