package models

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type ListItem struct {
	Id        uint      `gorm:"column:id;primary_key;auto_increment" json:"id"`
	DueDate   int64     `gorm:"column:due_date" json:"due_date"`
	Item      string    `gorm:"column:item" json:"item"`
	Done      bool      `gorm:"column:done" json:"done"`
	Over      bool      `gorm:"column:over" json:"over"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type ResponseItem struct {
	Id        uint   `gorm:"column:id;primary_key;auto_increment" json:"id"`
	DueDate   int64  `gorm:"column:due_date" json:"date"`
	Item      string `gorm:"column:item" json:"item"`
	Done      bool   `gorm:"column:done" json:"done"`
	Over      bool   `gorm:"column:over" json:"over"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func DbConnect() *gorm.DB {
	DB, err := gorm.Open(sqlite.Open("db/TodoList.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		err := DB.AutoMigrate(&ListItem{})
		if err != nil {
			println("failed to migrate todo list ")
			panic(fmt.Sprintf("Failed to migrate dababase %s\n", err))
		} else {
			return DB
		}
	}
}
