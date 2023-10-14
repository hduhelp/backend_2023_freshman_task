package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"login-system/db_handle"
	"login-system/utils"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// 检查是否提供了 JWT
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 解析 JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 检查签名方法是否有效
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法")
			}
			return utils.SecretKey, nil
		})

		// 验证 JWT
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
			c.Abort()
			return
		}

		// 在JWT验证通过后，检查JWT是否在黑名单内
		isBlacklisted, err := db_handle.IsTokenBlacklisted(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部错误"})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token已无效"})
			c.Abort()
			return
		}

		user, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "无效token"})
			return
		}

		c.Set("user", user)
		// JWT 验证通过，且不在黑名单内，继续处理请求
		c.Next()
	}
}
