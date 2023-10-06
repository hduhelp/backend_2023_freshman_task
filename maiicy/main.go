package main

import (
	"github.com/gin-gonic/gin"
	"login-system/db_handle"
	"login-system/routes"
)

func main() {
	_ = db_handle.ConnectDatabase("data.db")

	r := gin.Default()
	routes.SetupRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
