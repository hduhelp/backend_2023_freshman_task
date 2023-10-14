package main

import (
	"whxxxxxxxxxx/conf"
	"whxxxxxxxxxx/routes"
)

func main() {
	conf.Init()
	router := routes.NewRouter()
	router.Run(":" + conf.HttpPort)
}
