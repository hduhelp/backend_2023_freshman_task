// Package post 自动绑定json格式的新增请求 添加到数据库中
package post

import (
	"HDUhelper_Todo/models"

	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func (c Controller) Post(ctx *gin.Context) {
	tem := models.ListItem{}
	db := models.DbConnect()
	err := ctx.BindJSON(&tem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to bind",
		})
		return
	}
	err = db.Create(&tem).Error
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "failed to creat",
		})
		return
	}
	//tem.CreatedAt = time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Done!",
	})
}
