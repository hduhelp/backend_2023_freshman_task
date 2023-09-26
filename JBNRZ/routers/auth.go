package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo/models"
)

type PostForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(ctx *gin.Context) {
	var LoginForm PostForm
	var ret gin.H
	if ctx.ShouldBind(&LoginForm) == nil {
		username := LoginForm.Username
		password := LoginForm.Password
		status, _ := models.CheckAuth(username, password)
		ret = gin.H{
			"status": status.Description,
			"code":   status.Code,
			"time":   time.Now(),
			"data":   nil,
		}
		if status != models.LoginSuccess {
			models.Logger.Warning(ret)
			ctx.JSONP(http.StatusUnauthorized, ret)
			return
		} else {
			token := models.CreateToken(username)
			ret["data"] = gin.H{
				"username": username,
				"token":    token,
			}
			models.Logger.Info(ret)
			ctx.Header("Set-cookie", fmt.Sprintf("todo-token=%s;", token))
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		ret = gin.H{
			"status": "ParameterError",
			"code":   -1,
			"time":   time.Now(),
			"data":   nil,
		}
		models.Logger.Warning(ret)
		ctx.JSONP(http.StatusOK, ret)
		return
	}
}

func Register(ctx *gin.Context) {
	var RegisterForm PostForm
	var ret gin.H
	if ctx.ShouldBind(&RegisterForm) == nil {
		username := RegisterForm.Username
		password := RegisterForm.Password
		status, _ := models.AddUser(username, password)
		ret = gin.H{
			"status": status.Description,
			"code":   status.Code,
			"time":   time.Now(),
			"data":   nil,
		}
		if status != models.RegisterSuccess {
			models.Logger.Warning(ret)
			ctx.JSONP(http.StatusBadRequest, ret)
			return
		} else {
			token := models.CreateToken(username)
			ret["data"] = gin.H{
				"username": username,
				"token":    token,
			}
			models.Logger.Info(ret)
			ctx.Header("Set-cookie", fmt.Sprintf("todo-token=%s;", token))
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		ret = gin.H{
			"status": "ParameterError",
			"code":   -1,
			"time":   time.Now(),
			"data":   nil,
		}
		models.Logger.Warning(ret)
		ctx.JSONP(http.StatusOK, ret)
		return
	}
}
