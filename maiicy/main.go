package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"login-system/config"
	"login-system/db_handle"
	"login-system/routes"
	"login-system/utils"
	"os"
	"strconv"
)

func main() {

	databasePath := "data.db"
	configFilename := "config.toml"

	err := db_handle.ConnectDatabase(databasePath)
	if err != nil {
		panic(err)
	}

	MyConfig, err := config.LoadConfig(configFilename)
	if err != nil {
		fmt.Printf("无法读取配置文件：%v\n", err)
		os.Exit(1)
	}
	utils.SecretKey = []byte(MyConfig.SecretKey)

	r := gin.Default()
	routes.SetupRoutes(r)

	serverAddress := MyConfig.ServerIP + ":" + strconv.Itoa(MyConfig.ServerPort)
	err = r.Run(serverAddress)
	if err != nil {
		panic(err)
	}
}
