package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

type TODO struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Todos []TODO      `json:"todos"`
}

var todos []TODO
var nextID = 1

// 保存到文件
func saveFile() {
	file, err := os.Create("todos.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(todos)
	if err != nil {
		panic(err)
	}
}

// 读取文件
func loadFile() {
	file, err := os.Open("todos.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	if err != nil {
		panic(err)
	}
}

// 查询id
func findID(id int) int {
	for i, todo := range todos {
		if todo.ID == id {
			return i
		}
	}
	return -1
}

func main() {
	loadFile()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "%v", "Hello World")
	})

	// 添加todo
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		todo.ID = nextID
		nextID++
		todos = append(todos, todo)
		c.JSON(200, gin.H{"status": "ok"})
		saveFile()
	})

	// 删除todo
	r.DELETE("/todo/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		index := findID(id)
		if index != -1 {
			todos = append(todos[:index], todos[index+1:]...)
			c.JSON(200, gin.H{"status": "ok"})
		}
		saveFile()
	})

	// 修改todo
	r.PUT("/todo/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		index := findID(id)
		if index != -1 {
			var todo TODO
			c.BindJSON(&todo)
			todo.ID = id
			todos[index] = todo
			c.JSON(200, gin.H{"status": "ok"})
			saveFile()
		}
	})

	// 获取todo
	r.GET("/todo", func(c *gin.Context) {
		res := Response{
			Data:  struct{}{},
			Todos: todos,
		}
		c.JSON(200, res)
	})

	// 查询todo
	r.GET("/todo/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		index := findID(id)
		if index >= 0 && index < len(todos) {
			c.JSON(200, todos[index])
		}
	})

	r.Run(":8080")
}
