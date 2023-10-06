package mysql

import (
	"TODOlist/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitMysql() error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/hduhelp?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.TODO{})
	if err != nil {
		fmt.Println(err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
