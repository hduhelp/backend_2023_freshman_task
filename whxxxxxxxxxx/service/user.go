package service

import (
	"whxxxxxxxxxx/model"
	"whxxxxxxxxxx/pkg/utils"
	"whxxxxxxxxxx/serializer"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=12" msg:"用户名必须填写，且长度在3-12之间"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20" msg:"密码必须填写，且长度在6-15之间"`
}

func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count == 1 {
		return serializer.Response{
			Status: 40001,
			Msg:    "用户名已存在",
		}
	}
	user.UserName = service.UserName
	//对密码进行加密
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    "密码加密失败",
			Error:  err.Error(),
		}

	}
	//创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "注册时数据库错误",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "注册成功",
	}
}

func (service *UserService) Login() serializer.Response {
	var user model.User
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: 400,
				Msg:    "用户名不存在,请检查用户名或先注册",
				Error:  err.Error(),
			}
		}
		return serializer.Response{
			Status: 500,
			Msg:    "数据库错误",
			Error:  err.Error(),
		}
	}
	if !user.CheckPassword(service.Password) {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}
	//生成一个token用以前端交互
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "token生成失败",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: 200,
		//Data:   token,
		Data: serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:  "登录成功",
	}
}
