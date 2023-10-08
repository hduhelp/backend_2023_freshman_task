package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Todo struct {
	UUID    string `json:"uuid"`    // 唯一标识符
	Content string `json:"content"` // 内容
	Done    bool   `json:"done"`    // 是否完成
}

var todos []Todo
var filename string

func main() {
	filename, _ = filepath.Abs("todos.json")
	loadDataFromFile() // 从文件加载数据

	r := gin.Default()
	r.POST("/todo", createTodo)
	r.DELETE("/todo/:uuid", deleteTodo)
	r.PUT("/todo/:uuid", updateTodo)
	r.GET("/todos", getTodos)
	r.GET("/todo/:uuid", getTodo)
	r.Run(":8080")
}

func loadDataFromFile() {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file) // 读取文件内容
	if err != nil {
		log.Fatal(err)
	}

	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &todos) // 解析 JSON 数据到 todos 切片中
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveDataToFile() {
	bytes, err := json.Marshal(todos) // 将 todos 切片转换为 JSON 数据
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, bytes, 0666) // 将 JSON 数据写入文件
	if err != nil {
		log.Fatal(err)
	}
}

func createTodo(c *gin.Context) {
	var todo Todo
	err := c.BindJSON(&todo) // 解析请求中的 JSON 数据到 todo 结构体中
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	todo.UUID = uuid.New().String() // 生成唯一标识符

	todos = append(todos, todo) // 添加到 todos 切片中

	saveDataToFile() // 保存数据到文件

	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	uuid := c.Param("uuid")
	for i, todo := range todos {
		if todo.UUID == uuid {
			todos = append(todos[:i], todos[i+1:]...) // 从 todos 切片中删除指定元素

			saveDataToFile() // 保存数据到文件

			c.Status(http.StatusOK)
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func updateTodo(c *gin.Context) {
	uuid := c.Param("uuid")
	var updatedTodo Todo
	err := c.BindJSON(&updatedTodo) // 解析请求中的 JSON 数据到 updatedTodo 结构体中
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for i, todo := range todos {
		if todo.UUID == uuid {
			todos[i] = updatedTodo // 更新 todos 切片中的元素

			saveDataToFile() // 保存数据到文件

			c.Status(http.StatusOK)
			return
		}
	}
	c.Status(http.StatusNotFound)
}

func getTodos(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.DefaultQuery("page", "1")) // 获取页码参数，默认为第一页
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10")) // 获取每页条目数参数，默认为 10 条
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(todos) { // 如果起始索引超过 todos 切片长度，则返回空切片
		c.JSON(http.StatusOK, []Todo{})
		return
	}
	if endIndex > len(todos) {
		endIndex = len(todos) // 如果结束索引超过 todos 切片长度，则将结束索引设为 todos 切片的最后一个元素的索引加一
	}

	c.JSON(http.StatusOK, todos[startIndex:endIndex]) // 返回指定范围内的 todos 切片
}

func getTodo(c *gin.Context) {
	uuid := c.Param("uuid")
	for _, todo := range todos {
		if todo.UUID == uuid {
			c.JSON(http.StatusOK, todo) // 返回指定 UUID 的 todo
			return
		}
	}
	c.Status(http.StatusNotFound)
}
