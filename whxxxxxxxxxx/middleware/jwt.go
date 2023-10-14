package middleware

import (
	"time"
	"whxxxxxxxxxx/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		//var data interface{}
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = 403 //token无效
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 401
			}
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    "token验证失败",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
