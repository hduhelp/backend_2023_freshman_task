package api

import (
	"whxxxxxxxxxx/pkg/utils"
	"whxxxxxxxxxx/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func CreateTask(c *gin.Context) {
	var createTask service.CreateTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&createTask); err == nil {
		res := createTask.Create(claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func GetOneTask(c *gin.Context) {
	var getOneTask service.GetOneTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&getOneTask); err == nil {
		res := getOneTask.GetOne(c.Param("id"), claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func GetAllTask(c *gin.Context) {
	var getAllTask service.GetAllTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&getAllTask); err == nil {
		res := getAllTask.GetAll(claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func UpdateTask(c *gin.Context) {
	var updateTask service.UpdateTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&updateTask); err == nil {
		res := updateTask.Update(c.Param("id"), claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func SearchTask(c *gin.Context) {
	var searchTask service.SearchTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&searchTask); err == nil {
		res := searchTask.Search(claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}

func DeleteTask(c *gin.Context) {
	var deleteTask service.DeleteTaskService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&deleteTask); err == nil {
		res := deleteTask.Delete(c.Param("id"), claim.Id)
		c.JSON(200, res)

	} else {
		logging.Error(err)
		c.JSON(400, err)
	}
}
