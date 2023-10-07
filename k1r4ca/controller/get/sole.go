// Package getsole 获取指定日期的ToDo内容 传入指定日期0秒UNIX时间戳
package get

import (
	"HDUhelper_Todo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type SoleController struct {
}

func (c SoleController) GetSole(ctx *gin.Context) {

	dateStr := ctx.Param("date")

	//将UNIX时间戳转换为go time
	dateInt, _ := strconv.ParseInt(dateStr, 10, 64)
	date := time.Unix(dateInt, 0)
	db, err := gorm.Open(sqlite.Open("TodoList.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}

	Item := make([]*models.ListItem, 0)
	UpItem := make([]*models.ListItem, 0)

	//更新逾期item信息
	t := time.Now()
	err = db.Where("due_date <= ?", t).Find(&UpItem).Error
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

	//查询指定日期数据并以截止日期升序排序
	dd, _ := time.ParseDuration("24h")
	dbWhere := db.Where("created_at >= ?", date).Where("created_at <= ?", date.Add(dd))
	err = dbWhere.Order("due_date asc").Find(&Item).Error //太长了截成两句
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to find",
		})
		return
	}

	//用ResponseItem结构体框架下的Items进行时间转换 以传输unix时间戳
	Items := make([]*models.ResponseItem, len(Item))
	for i := range Item {
		Items[i] = &models.ResponseItem{
			Id:        Item[i].Id,
			DueDate:   Item[i].DueDate,
			Item:      Item[i].Item,
			Done:      Item[i].Done,
			Over:      Item[i].Over,
			CreatedAt: Item[i].CreatedAt.Unix(),
			UpdatedAt: Item[i].UpdatedAt.Unix(),
		}
	}

	//以json格式返回
	ctx.JSON(http.StatusOK, gin.H{
		"list": Items,
	})

}
