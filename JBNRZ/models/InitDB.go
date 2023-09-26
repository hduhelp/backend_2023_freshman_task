package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(Env.GetString("db")))
	if err != nil {
		Logger.Fatalln(err)
	}
	if err = db.AutoMigrate(&User{}, &Todo{}); err != nil {
		Logger.Fatalln(err)
	}
}
