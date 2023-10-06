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
