package middlewave

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo/models"
)

func Admin(ctx *gin.Context) {
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
		return
	}
	status, user := models.GetUserByName(username)
	ret := gin.H{
		"status": models.DatabaseError,
		"code":   models.DatabaseError.Code,
		"time":   time.Now(),
		"data":   nil,
	}
	if status != models.GetUserSuccess {
		models.Logger.Warning(ret)
		ctx.Abort()
		ctx.JSONP(http.StatusOK, ret)
		return
	}
	if user.Role {
		ctx.Set("Admin", username)
		ctx.Next()
	} else {
		ctx.Abort()
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
			"code":   401,
			"data":   nil,
		})
	}
}
