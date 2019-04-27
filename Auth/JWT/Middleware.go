package JWT

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Middleware gin中间件，检查token
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "无权限访问！",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken 解析token包含的信息，判断是否过期，有效性等
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "授权已过期！",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}
