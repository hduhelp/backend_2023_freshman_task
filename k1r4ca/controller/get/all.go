// Package getall 获取所有ToDo内容
package get

import (
	"HDUhelper_Todo/models"
	ulits "HDUhelper_Todo/utils"

	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AllController struct {
}

func (c AllController) GetAll(ctx *gin.Context) {

	Item := make([]*models.ListItem, 0)
	UpItem := make([]*models.ListItem, 0)

	db := models.DbConnect()

	//更新逾期item信息
	t := time.Now()
	err := db.Where("due_date <= ?", t).Find(&UpItem).Error
	for i := range UpItem {
		UpItem[i].Over = true
		UpItem[i].UpdatedAt = time.Now()
		err = db.Where("id = ?", UpItem[i].Id).Updates(&UpItem[i]).Error
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "failed to update To-do item",
		})
		return
	}

	err = db.Order("due_date asc").Find(&Item).Error //查询数据并以截止日期升序排序
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "failed to find",
		})
		return
	}

	//以json格式传输
	ctx.JSON(http.StatusOK, gin.H{
		"list": ulits.Response(Item), //时间转换 传输unix时间戳
	})
}
