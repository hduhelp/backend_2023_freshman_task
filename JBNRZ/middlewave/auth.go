package middlewave

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/models"
)

func Auth(ctx *gin.Context) {
	token, err := ctx.Cookie("todo-token")
	if err != nil {
		ctx.Abort()
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
			"code":   401,
			"data":   nil,
		})
	}
	var exists bool
	var username string
	models.Tokens.Range(func(key, value any) bool {
		if value.(string) == token {
			exists = true
			username = key.(string)
			return false
		}
		return true
	})
	if exists != true {
		ctx.Abort()
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
			"code":   401,
			"data":   nil,
		})
	}
	ctx.Set("Username", username)
	ctx.Next()
}
