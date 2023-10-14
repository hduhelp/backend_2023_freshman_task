package main

import (
	"TODOlist/controller"
	"TODOlist/dao/mysql"
	"TODOlist/middlewares/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//gorm
	err := mysql.InitMysql()
	if err != nil {
		fmt.Println(err)
	}

	//router
	r := gin.Default()
	//front
	r.Static("/static", "./static")
	r.LoadHTMLFiles("./view/index.html", "./view/todo.html")
	r.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	//user
	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)
	r.GET("/menu", jwt.MiddleWareJWT, func(c *gin.Context) {
		c.HTML(http.StatusOK, "todo.html", nil)
	})
	//todo
	todo := r.Group("/todo", jwt.MiddleWareJWT)
	{
		todo.GET("", controller.GetAllToDO)
		todo.POST("", controller.AddToDo)
		todo.DELETE("/:id", controller.VerifyPermission, controller.DeleteToDo)
		todo.PUT("/:id", controller.VerifyPermission, controller.UpdateToDo)
		todo.GET("/:id", controller.VerifyPermission, controller.GetTodo)
	}
	_ = r.Run(":8080")

}
