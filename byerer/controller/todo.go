package controller

import (
	"TODOlist/dao/mysql"
	"TODOlist/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func AddToDo(c *gin.Context) {
	var todo models.TODO
	_ = c.BindJSON(&todo)
	userID, ok := c.Get("userID")
	if !ok {
		//handle error
	}
	todo.UserID = userID.(int64)
	todo.ID = generateTodoID(todo.UserID)
	fmt.Println(todo)
	mysql.DB.Create(todo)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"todo":   todo,
	})
}

func DeleteToDo(c *gin.Context) {
	index := c.Param("id")
	result := mysql.DB.Where("id = ?", index).Delete(&models.TODO{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "record does not exist",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func UpdateToDo(c *gin.Context) {
	index := c.Param("id")
	var todo models.TODO
	_ = c.BindJSON(&todo)
	userID, ok := c.Get("userID")
	if !ok {

	}
	todo.UserID = userID.(int64)
	todo.ID = generateTodoID(todo.UserID)
	mysql.DB.Model(&models.TODO{}).Where("id = ?", index).Updates(todo)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func GetAllToDO(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {

	}
	var todos []models.TODO
	mysql.DB.Where("userID = ?", userID).Find(&todos)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"todos":  todos,
	})
}

func GetTodo(c *gin.Context) {
	index := c.Param("id")
	var todo models.TODO
	todo.ID = index
	mysql.DB.Where("ID = ?", todo.ID).Find(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"todo":   todo,
	})
}

// VerifyPermission 验证是否越权 写成中间件
func VerifyPermission(c *gin.Context) {
	userID, _ := c.Get("userID")
	if strconv.FormatInt(userID.(int64), 10) != c.Param("id")[:4] {
		c.JSON(http.StatusOK, gin.H{
			"message": "what are you fucking doing?",
		})
		c.Abort()
	}
	c.Next()
}

func generateTodoID(userID int64) string {
	Timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	todoID := strconv.FormatInt(userID, 10) + Timestamp
	return todoID
}
