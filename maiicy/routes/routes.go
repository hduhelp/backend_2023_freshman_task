package routes

import (
	"github.com/gin-gonic/gin"
	"login-system/handlers"
	"login-system/middlewares"
)

func SetupRoutes(r *gin.Engine) {
	// 用户资源
	userRouter := r.Group("/api/users")
	userRouter.POST("/register", handlers.RegisterHandler)
	userRouter.POST("/login", handlers.LoginHandler)
	userRouter.DELETE("/logout", middlewares.AuthMiddleware(), handlers.LogoutHandler)

	// 待办事项资源
	todoRouter := r.Group("/api/todos").Use(middlewares.AuthMiddleware())
	todoRouter.POST("/create", handlers.TodoAddHandler)          // 创建待办事项
	todoRouter.PUT("/:id", handlers.TodoUpdateHandler)           // 更新待办事项
	todoRouter.DELETE("/:id", handlers.TodoDelHandler)           // 删除待办事项
	todoRouter.GET("/:id", handlers.GetIDTodoHandler)            // 获取指定ID的待办事项
	todoRouter.GET("/before/:date", handlers.GetDateTodoHandler) // 获取指定日期前的待办事项
	todoRouter.GET("/all", handlers.GetAllTodoHandler)           // 获取所有待办事项
}
