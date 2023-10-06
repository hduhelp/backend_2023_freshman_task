package main

import (
	"todo/models"
	"todo/routers"
)

func main() {
	models.InitLogger()
	models.InitConfig()
	models.InitDB()
	models.InitAdmin()
	models.InitCron().Start()
	router := routers.InitRouters()
	if err := router.Run(":" + models.Env.GetString("server.port")); err != nil {
		models.Logger.Fatalln(err)
	}
}
