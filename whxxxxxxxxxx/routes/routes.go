package routes

import (
	"whxxxxxxxxxx/api"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/user/v1")
	{
		//用户操作
		v1.POST("/user/register", api.UserRegister)
		v1.POST("/user/login", api.UserLogin)
	}
	return r
}
