package ulits

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"

	"HDUhelper_Todo/models"
)

func IssueToken(username string) (string, error) {
	// 从配置文件中读取JWT签名密钥
	secret := viper.GetString("jwt.secret")

	// 创建JWT声明信息
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)), // 设置JWT过期时间为6小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "K1r4Ca",
		},
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用JWT签名密钥进行签名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头部信息
		authHeader := c.GetHeader("Authorization")

		// 检查Authorization头部是否为空
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You Need PERMISSION! ",
			})
			c.Abort() // 终止请求处理链
			return
		}

		// 检查Authorization头部是否以"Bearer "开头
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		// 解析JWT令牌
		tokenString := headerParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 从配置文件中读取JWT签名密钥
			secret := viper.GetString("jwt.secret")
			// 指定JWT签名密钥
			return []byte(secret), nil
		})

		// 检查JWT解析是否出错
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// 检查JWT是否有效
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token, Update it! ",
			})
			c.Abort()
			return
		}

		// 将解析后的JWT令牌存储在上下文中，以便其他处理程序使用
		c.Set("token", token)

		c.Next()
	}
}
