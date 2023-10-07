package main

import (
	"HDUhelper_Todo/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// 读取配置文件
	viper.SetConfigFile("configs/configs.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	r := gin.Default()

	routers.Init(r)

	_ = r.Run(":1070")

}
