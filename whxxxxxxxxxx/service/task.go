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

type UpdateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` //0未完成，1完成
}

type GetOneTaskService struct {
}

type DeleteTaskService struct {
}

type SearchTaskService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type GetAllTaskService struct {
	//分页功能
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
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

func (service *GetOneTaskService) GetOne(tid string) serializer.Response {
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "获取对应task失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "获取对应task成功",
		Data:   serializer.BuildTask(task),
	}
}

func (service *GetAllTaskService) GetAll(uid uint) serializer.Response {
	var tasks []model.Task
	var count int64
	count = 0
	//分页功能
	if service.PageNum == 0 {
		service.PageNum = 1
	}
	if service.PageSize == 0 {
		service.PageSize = 10
	}
	model.DB.Model(&model.Task{}).Preload("User").Where("uid=?", uid).Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks)
	/* return serializer.Response{
		Status: 200,
		Msg:    "获取所有task成功",
		Data:   serializer.BuildTasks(tasks),
	} */
	//添加分页功能
	//print(tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

func (service *UpdateTaskService) Update(tid string) serializer.Response {
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "更新时获取对应task失败",
		}
	}
	task.Title = service.Title
	task.Content = service.Content
	task.Status = service.Status
	err = model.DB.Save(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "更新对应task失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "更新对应task成功",
		Data:   serializer.BuildTask(task),
	}
}

func (service *SearchTaskService) Search(uid uint) serializer.Response {
	var tasks []model.Task
	var count int64
	count = 0
	if service.PageNum == 0 {
		service.PageNum = 1
	}
	if service.PageSize == 0 {
		service.PageSize = 10
	}
	model.DB.Model(&model.Task{}).Preload("User").Where("uid=?", uid).Where("title LIKE ? OR content LIKE ?", "%"+service.Info+"%", "%"+service.Info+"%").Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

func (service *DeleteTaskService) Delete(tid string) serializer.Response {
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "删除时获取对应task失败",
		}
	}
	err = model.DB.Delete(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "删除对应task失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "删除对应task成功",
	}
}
