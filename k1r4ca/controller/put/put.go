// Package put 更新ToDo内容
package put

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"HDUhelper_Todo/models"
)

type Controller struct {
}

func (c Controller) Put(ctx *gin.Context) {

	tem := models.ListItem{}

	err := ctx.BindJSON(&tem)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "failed to update",
		})
		return
	}

	db := models.DbConnect()
	tem.UpdatedAt = time.Now()
	err = db.Where("id = ?", tem.Id).Updates(&tem).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to update",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Done!",
	})
}
