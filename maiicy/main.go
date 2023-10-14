package main

import (
	"flag"
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
	// 读取配置文件
	configFilename := "config.toml"
	MyConfig, err := config.LoadConfig(configFilename)
	if err != nil {
		fmt.Printf("无法读取配置文件：%v\n", err)
		os.Exit(1)
	}
	utils.SecretKey = []byte(MyConfig.SecretKey)

	databasePath := MyConfig.DBPath

	// 读取cli参数
	dbPath := flag.String("dbPath", "null", "数据库路径")
	host := flag.String("host", "null", "服务器主机名")
	port := flag.Int("port", -1, "服务器端口号")

	flag.Parse()

	// 如果cli中有输入，则用cli参数，如果没有输入，则用配置文件内数据
	if *dbPath != "null" {
		databasePath = *dbPath
	}
	if *host != "null" {
		MyConfig.ServerIP = *host
	}
	if *port != -1 {
		MyConfig.ServerPort = *port
	}

	err = db_handle.ConnectDatabase(databasePath)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	serverAddress := MyConfig.ServerIP + ":" + strconv.Itoa(MyConfig.ServerPort)
	err = r.Run(serverAddress)
	if err != nil {
		panic(err)
	}
}
