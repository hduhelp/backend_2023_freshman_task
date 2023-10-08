package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TODO struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var todos []TODO

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {
	r := gin.Default()

	// 添加TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: err.Error()})
			return
		}
		todos = append(todos, todo)
		c.JSON(http.StatusOK, Response{Status: "ok"})
	})

	// 删除TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: err.Error()})
			return
		}
		if index < 0 || index >= len(todos) {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: "index out of range"})
			return
		}
		todos = append(todos[:index], todos[index+1:]...)
		c.JSON(http.StatusOK, Response{Status: "ok"})
	})

	// 修改TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: err.Error()})
			return
		}
		var todo TODO
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: err.Error()})
			return
		}
		todos[index] = todo
		c.JSON(http.StatusOK, Response{Status: "ok"})
	})

	// 获取TODO列表
	r.GET("/todo", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})

	// 获取TODO详情
	r.GET("/todo/:index", func(c *gin.Context) {
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: err.Error()})
			return
		}
		if index < 0 || index >= len(todos) {
			c.JSON(http.StatusBadRequest, Response{Status: "error", Data: "index out of range"})
			return
		}
		c.JSON(http.StatusOK, todos[index])
	})

	r.Run(":8080")
}
