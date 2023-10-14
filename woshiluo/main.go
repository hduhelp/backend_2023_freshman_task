package main

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"todolist/controller"
	"todolist/models"
	"todolist/utils"
)

func main() {
	if _, err := toml.DecodeFile("config.toml", &utils.ConfigData); err != nil {
		panic(err)
	}

	log.SetLevel(log.TraceLevel)
	models.ConnectDatabase(utils.ConfigData.DatabaseFile)

	log.Trace(utils.ConfigData)

	go utils.Monitor()
	r := gin.Default()

	r.GET("/todo", controller.ListTodo)
	r.POST("/todo", controller.NewTodo)
	r.GET("/todo/:id", controller.GetTodo)
	r.PUT("/todo/:id", controller.UpdateTodo)
	r.DELETE("/todo/:id", controller.DeleteTodo)

	r.POST("/user", controller.NewUser)
	r.PUT("/user/:id", controller.UpdateUser)

	r.POST("/token", controller.NewToken)
	r.DELETE("/token/:token", controller.DeleteToken)

	r.Run()
}
