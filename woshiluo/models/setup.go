//
// models.go
// Copyright (C) 2023 Woshiluo Luo <woshiluo.luo@outlook.com>
//
// Distributed under terms of the GNU AGPLv3+ license.
//

package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Todo struct {
	gorm.Model
	Title    string  `json:"title" binding:"required"`
	Location string  `json:"location"`
	DueDate  uint64  `json:"duedate"`
	UserID   uint    `json:"userid" binding:"required"`
	Done     bool    `json:"done"`
}

type Token struct {
	gorm.Model
	Token string `json:"token" binding:"required"`
	UserID uint  `json:"userid" binding:"required"`
}

var Db *gorm.DB

func ConnectDatabase(database_file string) {
	database, err := gorm.Open(sqlite.Open(database_file), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := database.AutoMigrate(&User{},&Todo{},&Token{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate dababase %s", err));
	}

	Db = database
}
