package models

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	Info     string
}

type Todo struct {
	TodoID    uint `gorm:"primaryKey;autoIncrement"`
	UserID    uint
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	DueDate   time.Time
}

type JWTBlacklist struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	Token  string    `gorm:"not null"`
	Expiry time.Time `gorm:"not null"`
}
