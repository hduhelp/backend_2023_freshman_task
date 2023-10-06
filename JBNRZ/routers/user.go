package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo/models"
)

var StatusUnauthorized = gin.H{
	"status": "Unauthorized",
	"code":   401,
	"time":   time.Now(),
	"data":   nil,
}

var StatusParameterError = gin.H{
	"status": "ParameterError",
	"code":   -1,
	"time":   time.Now(),
	"data":   nil,
}

type ChangePWDForm struct {
	OldPWD string `form:"old" binding:"required"`
	NewPWd string `form:"new" binding:"required"`
}

type AddTodoForm struct {
	ItemName string `form:"name" binding:"required"`
	Detail   string `form:"detail" binding:"required"`
	// + second
	EndTime int64 `form:"time" binding:"required"`
}

type ListTodoForm struct {
	From     *int `form:"from" binding:"required"`
	Quantity *int `form:"num" binding:"required"`
}

type DelTodoForm struct {
	ItemName string `form:"name" binding:"required"`
}

type GetTodoForm struct {
	ItemName string `form:"name" binding:"required"`
}

type ChangeTodoForm struct {
	ItemName string `form:"name" binding:"required"`
	// + second
	EndTime int64 `form:"time" binding:"required"`
}

type SetEmailForm struct {
	Email string `form:"email" binding:"required"`
}

func match(ctx *gin.Context) (bool, string, any, map[string]interface{}) {
	username := ctx.Param("username")
	token, ok := models.Tokens.Load(username)
	gcToken, err := ctx.Cookie("todo-token")
	if !ok || err != nil || token != gcToken {
		return false, "", "", StatusUnauthorized
	}
	return true, username, token, gin.H{}
}

func HomePage(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
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
			"Title":  username,
			"Logout": "/user/" + username + "/logout",
		})
		return
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func Reset(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var ChangePWD ChangePWDForm
		if err := ctx.ShouldBind(&ChangePWD); err == nil {
			status, user := models.ChangePWD(username, ChangePWD.OldPWD, ChangePWD.NewPWd)
			ret = gin.H{
				"status": status.Description,
				"code":   status.Code,
				"time":   time.Now(),
				"data":   nil,
			}
			if status != models.ChangePWDSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["data"] = gin.H{
					"username": user.Username,
					"token":    token,
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, StatusParameterError)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func Logout(ctx *gin.Context) {
	status, username, _, ret := match(ctx)
	if status {
		models.Tokens.Delete(username)
		ret = gin.H{
			"status": "success",
			"code":   0,
			"time":   time.Now(),
			"data": gin.H{
				"username": username,
				"token":    nil,
			},
		}
		models.Logger.Info(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusOK, ret)
		return
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func AddTodo(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var AddForm AddTodoForm
		ret = StatusParameterError
		if err := ctx.ShouldBind(&AddForm); err == nil {
			status, todo := models.AddTodo(AddForm.ItemName, AddForm.Detail, AddForm.EndTime, username)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.AddTodoSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"todo": gin.H{
						"name":   todo.ItemName,
						"end":    todo.EndTime,
						"detail": todo.Detail,
					},
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func ListTodo(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var ListForm ListTodoForm
		ret = StatusParameterError

		if err := ctx.ShouldBindQuery(&ListForm); err == nil {
			status, todos := models.ListTodo(username, *ListForm.From, *ListForm.Quantity)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.ListTodoSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"todos":    todos,
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func DelTodo(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var DelForm DelTodoForm
		ret = StatusParameterError
		if err := ctx.ShouldBind(&DelForm); err == nil {
			status, _ := models.DelTodo(DelForm.ItemName, username)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.DelTodoSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["status"] = status.Description
				ret["code"] = status.Code
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"deleted":  DelForm.ItemName,
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func GetTodo(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var GetForm GetTodoForm
		ret = StatusParameterError
		if err := ctx.ShouldBindQuery(&GetForm); err == nil {
			status, todo := models.GetTodo(GetForm.ItemName, username)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.ListTodoSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["status"] = status.Description
				ret["code"] = status.Code
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"todo": gin.H{
						"name":   todo.ItemName,
						"detail": todo.Detail,
						"end":    todo.EndTime,
					},
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func ChangeTodo(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var ChangeForm ChangeTodoForm
		ret = StatusParameterError
		if err := ctx.ShouldBind(&ChangeForm); err == nil {
			status, todo := models.ChangeTodo(ChangeForm.ItemName, username, ChangeForm.EndTime)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.ChangeTodoSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["status"] = status.Description
				ret["code"] = status.Code
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"update": gin.H{
						"name":   todo.ItemName,
						"detail": todo.Detail,
						"end":    todo.EndTime,
					},
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}

func SetEmail(ctx *gin.Context) {
	status, username, token, ret := match(ctx)
	if status {
		var SetForm SetEmailForm
		ret = StatusParameterError
		if err := ctx.ShouldBind(&SetForm); err == nil {
			status, user := models.SetEmail(username, SetForm.Email)
			ret["status"] = status.Description
			ret["code"] = status.Code
			ret["data"] = nil
			if status != models.SetEmailSuccess {
				models.Logger.Warning(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			} else {
				ret["status"] = status.Description
				ret["code"] = status.Code
				ret["data"] = gin.H{
					"username": username,
					"token":    token,
					"email":    user.Email,
				}
				models.Logger.Info(ret)
				ctx.JSONP(http.StatusOK, ret)
				return
			}
		} else {
			models.Logger.Warning(err)
			ctx.JSONP(http.StatusOK, ret)
			return
		}
	} else {
		models.Logger.Warning(ret)
		ctx.Redirect(http.StatusFound, "/login")
		//ctx.JSONP(http.StatusUnauthorized, ret)
		return
	}
}
