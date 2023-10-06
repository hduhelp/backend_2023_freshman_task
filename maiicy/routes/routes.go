package routes

import (
	"github.com/gin-gonic/gin"
	"login-system/handlers"
	"login-system/middlewares"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/user/register", handlers.RegisterHandler)
	r.POST("/api/user/login", handlers.LoginHandler)
	r.POST("/api/user/logout", middlewares.AuthMiddleware(), handlers.LogoutHandler)

	r.POST("/api/todo/add", middlewares.AuthMiddleware(), handlers.TodoAddHandler)
	r.POST("/api/todo/delete", middlewares.AuthMiddleware(), handlers.TodoDelHandler)
	r.POST("/api/todo/update", middlewares.AuthMiddleware(), handlers.TodoUpdateHandler)
	r.GET("/api/todo/:date", middlewares.AuthMiddleware(), handlers.GetDateTodoHandler)
	r.GET("/api/todo/all", middlewares.AuthMiddleware(), handlers.GetAllTodoHandler)
}
