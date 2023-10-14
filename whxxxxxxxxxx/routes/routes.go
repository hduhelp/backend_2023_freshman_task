package routes

import (
	"fmt"
	"whxxxxxxxxxx/api"
	"whxxxxxxxxxx/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/user/v2")
	{
		//用户操作
		fmt.Println("路由1")
		v1.POST("/register", api.UserRegister)
		v1.POST("/login", api.UserLogin)

	}
	authed := r.Group("api/task/v2")
	authed.Use(middleware.JWT())
	{
		authed.POST("/create", api.CreateTask)
		authed.GET("/getone/:id", api.GetOneTask)
		authed.GET("/getall", api.GetAllTask)
		authed.PUT("/update/:id", api.UpdateTask)
		authed.POST("/search", api.SearchTask)
		authed.DELETE("/delete/:id", api.DeleteTask)
	}
	return r
}
