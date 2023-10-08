package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

type Data struct {
	Users []User `json:"users"`
	Todos []TODO `json:"todos"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TODO struct {
	Content  string `json:"content"`
	Done     bool   `json:"done"`
	Deadline string `json:"deadline"`
}

var data *Data

func loadData() *Data {
	data := &Data{}
	os.Create("data.json")
	b, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("新建文件")
		return data
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		fmt.Println("错误:", err)
	}
	return data
}

func saveData(data *Data) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("错误:", err)
		return
	}
	err = os.WriteFile("data.json", b, 0644)
	if err != nil {
		fmt.Println("错误:", err)
	}
}

func main() {
	r := gin.Default()

	// 登录接口
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 验证用户名和密码是否正确
		if username == "admin" && password == "123" {
			c.JSON(http.StatusOK, gin.H{
				"code": "0",
				"msg":  "登录成功",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "用户名或密码错误",
			})
		}
	})

	// 新增 todo
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		data.Todos = append(data.Todos, todo)

		fmt.Println(data.Todos)

		c.JSON(200, gin.H{"status": "ok"})
	})

	// 删除 todo
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		data.Todos = append(data.Todos[:index], data.Todos[index+1:]...)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 更新 todo
	r.PUT("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		data.Todos[index] = todo
		c.JSON(200, gin.H{"status": "ok"})
	})

	//获取 todo
	r.GET("/todo", func(c *gin.Context) {
		c.JSON(200, data.Todos)
	})

	// 查询 todo
	r.GET("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		c.JSON(200, data.Todos[index])

	})

	// 写入数据
	data = loadData()
	defer saveData(data)

	r.Run(":8080")
}
