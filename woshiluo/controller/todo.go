//
// todo.go
// Copyright (C) 2023 Woshiluo Luo <woshiluo.luo@outlook.com>
//
// Distributed under terms of the GNU AGPLv3+ license.
//

package controller

import (
	_ "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"todolist/models"
)

type ListTodoData struct {
	Token string `json:"token" binding:"required"`
}

func ListTodo(c *gin.Context) {
	var data ListTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}
	var todos []models.Todo
	models.Db.Model(models.Todo{UserID: user.ID}).Find(&todos)
	c.JSON(http.StatusOK, todos)
}

type GetTodoData struct {
	Token string `json:"token" binding:"required"`
}

func GetTodo(c *gin.Context) {
	var id = c.Param("id")
	var todo models.Todo
	var data GetTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}

	if err := models.Db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "Todo Not Found"})
		return
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

type NewTodoData struct {
	Token string      `json:"token" binding:"required"`
	Todo  models.Todo `json:"todo" binding:"required"`
}

func NewTodo(c *gin.Context) {
	var data NewTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}

	if data.Todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Forbidden"})
		return
	}

	if err := models.Db.Create(&data.Todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.Todo)
}

type UpdateTodoData struct {
	Token string      `json:"token" binding:"required"`
	Todo  models.Todo `json:"todo" binding:"required"`
}

func UpdateTodo(c *gin.Context) {
	var origin_todo models.Todo
	var id = c.Param("id")
	var data UpdateTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}

	if err := models.Db.First(&origin_todo, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Todo Not Found"})
		return
	}

	if user.ID != origin_todo.UserID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Forbidden"})
		return
	}

	// NOTE When update with struct, GORM will only update non-zero fields, you might want to use map to update attributes or use Select to specify fields to update
	// Source: https://stackoverflow.com/questions/56653423/gorm-doesnt-update-boolean-field-to-false

	if data.Todo.Done == false {
		if err := models.Db.Model(&origin_todo).Updates(map[string]interface{}{"done": false}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	if err := models.Db.Model(&origin_todo).Updates(data.Todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, origin_todo)
}

type DeleteTodoData struct {
	Token string      `json:"token" binding:"required"`
	Todo  models.Todo `json:"todo" binding:"required"`
}

func DeleteTodo(c *gin.Context) {
	var todo models.Todo
	var id = c.Param("id")
	var data DeleteTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}

	if err := models.Db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if user.ID != todo.UserID {
		c.JSON(http.StatusForbidden, gin.H{"err": "Forbidden"})
		return
	}

	if err := models.Db.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}
