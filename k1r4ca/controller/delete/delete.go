// Package delete_todo 删除指定id的ToDo内容
package delete

import (
	"HDUhelper_Todo/controller/delete/DelTx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
}

func (c Controller) Delete(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr) //利用int分析将idStr转换为int64

	err := DelTx.DelTx(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to delete",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Done!",
	})

}
