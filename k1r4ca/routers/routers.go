package routers

import (
	"HDUhelper_Todo/controller/delete"
	"HDUhelper_Todo/controller/get"
	"HDUhelper_Todo/controller/login"
	"HDUhelper_Todo/controller/post"
	"HDUhelper_Todo/controller/put"

	ulits "HDUhelper_Todo/utils"

	"github.com/gin-gonic/gin"
)

func Init(g *gin.Engine) {

	rt := g.Group("/list", ulits.AuthMiddleware())
	{
		// POST 路由
		rt.POST("/", post.Controller{}.Post)

		// DELETE 路由
		rt.DELETE("/:id", delete.Controller{}.Delete)

		// PUT 路由
		rt.PUT("/", put.Controller{}.Put)

		// GET 路由
		rt.GET("/", get.AllController{}.GetAll)
		rt.GET("/:date", get.SoleController{}.GetSole)

	}

	g.GET("/login/:username", login.IssueJWT)

}
