package routers

import (
	"github.com/gin-gonic/gin"
	"todo/middlewave"
)

func InitRouters() (router *gin.Engine) {
	router = gin.Default()

	defaultRouters := router.Group("/")
	defaultRouters.GET("/", RootPage)
	defaultRouters.GET("/login", LoginPage)
	defaultRouters.GET("/register", RegisterPage)
	defaultRouters.POST("/login", Login)
	defaultRouters.POST("/register", Register)

	userRouters := router.Group("/user", middlewave.Auth)
	userRouters.GET("/:username/home", HomePage)
	userRouters.POST("/:username/reset", Reset)
	userRouters.GET("/:username/logout", Logout)
	userRouters.POST("/:username/add", AddTodo)
	userRouters.GET("/:username/list", ListTodo)
	userRouters.GET("/:username/get", GetTodo)
	userRouters.POST("/:username/delete", DelTodo)
	userRouters.POST("/:username/update", ChangeTodo)
	userRouters.POST("/:username/email", SetEmail)

	adminRouters := router.Group("/admin", middlewave.Admin)
	adminRouters.GET("/list", ListAll)
	adminRouters.POST("/add", AdminAddTodo)
	adminRouters.POST("/delete", AdminDelTodo)
	return router
}
