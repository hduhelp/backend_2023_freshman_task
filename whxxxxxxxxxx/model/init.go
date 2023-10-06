package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DatabaseLink(constring string) {
	fmt.Printf("constring: %v\n", constring)
	var err error
	/* 	//后续可以写到配置文件中
	   	dsn := "root:123456@tcp(127.0.0.1:3306)/whxxxxxxxxxx?charset=utf8mb4&parseTime=True&loc=Local" */
	//链接数据库
	DB, err = gorm.Open(mysql.Open(constring), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败")
	} else {
		fmt.Println("数据库连接成功")
	}
	//记录所有日志
	DB.Logger.LogMode(logger.Info)
	if gin.Mode() == "release" {
		//不记录日志
		DB.Logger.LogMode(logger.Silent)
	}
	//根据结构体格式创建数据表
	//db.AutoMigrate(&todoModel{})
	migration()
}
