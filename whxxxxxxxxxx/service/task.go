package service

import (
	"time"
	"whxxxxxxxxxx/model"
	"whxxxxxxxxxx/serializer"
)

type CreateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` //0未完成，1完成
}

func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	code := 200
	model.DB.First(&user, id)
	task := model.Task{
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Status:    0,
		Content:   service.Content,
		StartTime: time.Now().Unix(),
		EndTime:   0,
	}
	err := model.DB.Create(&task).Error
	if err != nil {
		code = 500 //创建不成功
		return serializer.Response{
			Status: code,
			Msg:    "创建task失败",
		}

	}
	return serializer.Response{
		Status: code,
		Msg:    "创建task成功",
	}
}
