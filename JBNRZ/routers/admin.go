package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo/models"
)

type AdminAddTodoForm struct {
	ItemName string `form:"name" binding:"required"`
	Detail   string `form:"detail" binding:"required"`
	// + second
	EndTime  int64  `form:"time" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type AdminDelTodoForm struct {
	ItemName string `form:"name" binding:"required"`
	Username string `form:"username" binding:"required"`
}

func AdminHome(ctx *gin.Context) {
	var (
		username string
		ret      map[string]interface{}
	)
	token, ok := models.Tokens.Load("admin")
	gcToken, err := ctx.Cookie("todo-token")
	status := false
	if ok || err == nil || token == gcToken {
		status = true
		username = "admin"
	}
	if status {
		ret = gin.H{
			"status": "success",
			"code":   0,
			"time":   time.Now(),
			"data": gin.H{
				"username": username,
				"token":    token,
			},
		}
		models.Logger.Info(ret)
		//ctx.JSONP(http.StatusOK, ret)
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Title": username,
		})
		return
	} else {
		models.Logger.Warning(StatusUnauthorized)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, StatusUnauthorized)
		return
	}
}

func ListAll(ctx *gin.Context) {
	var ListForm ListTodoForm
	ret := StatusParameterError
	if ctx.ShouldBindQuery(&ListForm) == nil {
		status, todos := models.ListAll(*ListForm.From, *ListForm.Quantity)
		ret["status"] = status.Description
		ret["code"] = status.Code
		ret["data"] = nil
		if status != models.ListTodoSuccess {
			models.Logger.Warning(ret)
			ctx.JSONP(http.StatusOK, ret)
		} else {
			ret["data"] = gin.H{
				"username": "admin",
				"token":    "admin-token",
				"todos":    todos,
			}
			models.Logger.Info(ret)
			ctx.JSONP(http.StatusOK, ret)
		}
	} else {
		models.Logger.Warning(ret)
		ctx.JSONP(http.StatusOK, ret)
	}
}

func AdminAddTodo(ctx *gin.Context) {
	var AddForm AdminAddTodoForm
	ret := StatusParameterError
	if err := ctx.ShouldBind(&AddForm); err == nil {
		status, todo := models.AddTodo(AddForm.ItemName, AddForm.Detail, AddForm.EndTime, AddForm.Username)
		ret["status"] = status.Description
		ret["code"] = status.Code
		ret["data"] = nil
		if status != models.AddTodoSuccess {
			models.Logger.Warning(ret)
			ctx.JSONP(http.StatusOK, ret)
		} else {
			ret["data"] = gin.H{
				"username": "admin",
				"token":    "admin-token",
				"todo": gin.H{
					"name":   todo.ItemName,
					"end":    todo.EndTime,
					"detail": todo.Detail,
				},
			}
			models.Logger.Info(ret)
			ctx.JSONP(http.StatusOK, ret)
		}
	} else {
		models.Logger.Warning(err)
		ctx.JSONP(http.StatusOK, ret)
	}
}

func AdminDelTodo(ctx *gin.Context) {
	var DelForm AdminDelTodoForm
	ret := StatusParameterError
	if err := ctx.ShouldBind(&DelForm); err == nil {
		status, _ := models.DelTodo(DelForm.ItemName, DelForm.Username)
		ret["status"] = status.Description
		ret["code"] = status.Code
		ret["data"] = nil
		if status != models.DelTodoSuccess {
			models.Logger.Warning(ret)
			ctx.JSONP(http.StatusOK, ret)
		} else {
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = gin.H{
				"username": "admin",
				"token":    "admin-token",
				"deleted":  DelForm.ItemName,
			}
			models.Logger.Info(ret)
			ctx.JSONP(http.StatusOK, ret)
		}
	} else {
		models.Logger.Warning(err)
		ctx.JSONP(http.StatusOK, ret)
	}
}
