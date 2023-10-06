package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	User      User   `gorm:"foreignKey:Uid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index;not null"` //查询索引
	Status    int    `gorm:"default:0"`      //0未完成，1完成
	Content   string `gorm:"type:longtext"`
	StartTime int64  //list开始时间
	EndTime   int64  //list结束时间
}
