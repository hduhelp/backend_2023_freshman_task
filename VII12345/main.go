package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TODO struct {
	Number      string `json:"number"`
	Content     string `json:"content"`
	Done        string `json:"done"`
	Finish_time string `json:"finish_time"`
}

type PASS struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var todos []TODO

var passes []PASS

func reg(text PASS) (a int) {
	file, _ := os.OpenFile("userdata.txt", os.O_RDWR, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	file.Close()
	if len(text.Username) <= 7 {
		return 0
	}
	for i := 0; i < len(readline); i = i + 2 {
		if readline[i] == text.Username {
			return 1
		}
	}
	return 2
}

func log(user string, pas string) (a bool) {
	file, _ := os.OpenFile("userdata.txt", os.O_APPEND, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	file.Close()
	lenth := len(readline) / 2
	for i := 0; i < lenth; i++ {
		if readline[i*2] == user && readline[i*2+1] == pas {
			return true
		}
	}
	return false
}

func check(user string) (a bool) {
	file, _ := os.OpenFile("userdata.txt", os.O_APPEND, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	file.Close()
	lenth := len(readline) / 2
	for i := 0; i < lenth; i++ {
		if readline[i*2] == user {
			return true
		}
	}
	return false
}

//读取过去数据
func readDATA(user string) {
	name := user + ".txt"
	file, _ := os.OpenFile(name, os.O_RDWR, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	lenth := len(readline) / 4
	var todo TODO
	for i := 0; i < lenth; i++ {
		todo.Number, todo.Content, todo.Done, todo.Finish_time = readline[i*4], readline[i*4+1], readline[i*4+2], readline[i*4+3]
		todos = append(todos, todo)
	}
	file.Close()

	file_p, _ := os.OpenFile("userdata.txt", os.O_RDWR, 0666)
	read_p, _ := io.ReadAll(file_p)
	readline_p := strings.Split(string(read_p), "\n")
	lenth_p := len(readline_p) / 2
	var pass PASS
	for i := 0; i < lenth_p; i++ {
		pass.Username, pass.Password = readline_p[i*2], readline_p[i*2+1]
		passes = append(passes, pass)
	}
	file.Close()
}

func add(text TODO, user string) {
	name := user + ".txt"
	file, _ := os.OpenFile(name, os.O_APPEND, 0666)
	file.WriteString(text.Number)
	file.WriteString("\n")
	file.WriteString(text.Content)
	file.WriteString("\n")
	file.WriteString(text.Done)
	file.WriteString("\n")
	file.WriteString(text.Finish_time)
	file.WriteString("\n")
	file.Close()
}

func del(ind int, user string) {
	name := user + ".txt"
	file, _ := os.OpenFile(name, os.O_RDWR, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	readline[ind*4+1], readline[ind*4+2], readline[ind*4+3] = `已被删除`, `已被删除`, `已被删除`
	file.Truncate(0)
	file.Seek(0, 0)
	for _, i := range readline {
		file.WriteString(i)
		file.WriteString("\n")
	}
	file.Close()
}

func change(text TODO, ind int, user string) {
	name := user + ".txt"
	file, _ := os.OpenFile(name, os.O_RDWR, 0666)
	read, _ := io.ReadAll(file)
	readline := strings.Split(string(read), "\n")
	ind = ind - 1
	readline[ind*4+1], readline[ind*4+2], readline[ind*4+3] = text.Content, text.Done, text.Finish_time
	file.Truncate(0)
	file.Seek(0, 0)
	for _, i := range readline {
		file.WriteString(i)
		file.WriteString("\n")
	}
	file.Close()
}

func main() {
	r := gin.Default()

	//增加 TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		username := c.Query("username")
		if check(username) {
			readDATA(username)
			todo.Number = strconv.Itoa(len(todos) + 1)
			todos = append(todos, todo)
			fmt.Println(todos)
			add(todo, username)
			c.JSON(200, gin.H{"status": "成功增加TODO"})
			todos = todos[:0]
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}

	})

	//删除 TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		username := c.Query("username")
		readDATA(username)
		if check(username) {
			if index > len(todos) {
				c.JSON(200, gin.H{"status": "序号过大，您没有这么多TODO"})
			} else if index <= 0 {
				c.JSON(200, gin.H{"status": "序号需大于0"})
			} else {
				todos[index-1].Content = `已被删除`
				todos[index-1].Done = `已被删除`
				todos[index-1].Finish_time = `已被删除`
				del(index-1, username)
				c.JSON(200, gin.H{"status": "成功删除TODO"})
			}
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}
		todos = todos[:0]
	})

	//更新 TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		username := c.Query("username")
		readDATA(username)
		if check(username) {
			if todos[index-1].Content == "已被删除" {
				c.JSON(200, gin.H{"status": "该TODO已被删除"})
			} else if index > len(todos) {
				c.JSON(200, gin.H{"status": "序号过大,您没有这么多TODO"})
			} else if index <= 0 {
				c.JSON(200, gin.H{"status": "序号需大于0"})
			} else {
				todos[index-1] = todo
				change(todo, index, username)
				c.JSON(200, gin.H{"status": "成功更新TODO"})
			}
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}
		todos = todos[:0]
	})

	//列出 TODO
	r.GET("/todo", func(c *gin.Context) {
		page := c.Query("page")
		intpage, _ := strconv.Atoi(page)
		username := c.Query("username")
		readDATA(username)
		if check(username) {
			if intpage <= 0 {
				c.JSON(200, gin.H{"status": "页码需大于等于1"})
			} else {
				if (intpage-1)*20 > len(todos) {
					c.JSON(200, gin.H{"status": "页码过大，您没有这么多TODO"})
				} else if len(todos) == 0 {
					c.JSON(200, gin.H{"status": "您还未添加TODO"})
				} else {
					c.JSON(200, gin.H{"status": "以下是该页的TODO内容："})
					if intpage*20 > len(todos) {
						c.JSON(200, todos[20*(intpage-1):])
					} else {
						c.JSON(200, todos[20*(intpage-1):20*intpage])
					}
				}
			}
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}
		todos = todos[:0]
	})

	//查询 TODO
	r.GET("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		username := c.Query("username")
		readDATA(username)
		if check(username) {
			if index > len(todos) {
				c.JSON(200, gin.H{"status": "序号过大,请重新查询"})
			} else if index <= 0 {
				c.JSON(200, gin.H{"status": "序号需大于0"})
			} else {
				c.JSON(200, gin.H{"status": "以下是您所查询的的TODO："})
				c.JSON(200, todos[index-1])
			}
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}
		todos = todos[:0]
	})

	//清空 TODO
	r.DELETE("/todo/delete", func(c *gin.Context) {
		username := c.Query("username")
		if check(username) {
			name := username + ".txt"
			file, _ := os.OpenFile(name, os.O_RDWR, 0666)
			file.Truncate(0)
			file.Seek(0, 0)
			file.Close()
			c.JSON(200, gin.H{"status": "已成功清空TODO"})
			todos = todos[:0]
		} else {
			c.JSON(200, gin.H{"status": "请先登录"})
		}
	})

	//注册
	r.POST("/register", func(c *gin.Context) {
		var pass PASS
		c.BindJSON(&pass)
		flag := reg(pass)
		if flag == 2 {
			file, _ := os.OpenFile("userdata.txt", os.O_APPEND, 0666)
			file.WriteString(pass.Username)
			file.WriteString("\n")
			file.WriteString(pass.Password)
			file.WriteString("\n")
			file.Close()
			passes = append(passes, pass)
			name := pass.Username + ".txt"
			os.Create(name)
			c.JSON(200, gin.H{"status": "注册成功"})
		} else if flag == 0 {
			c.JSON(200, gin.H{"status": "用户名应不小于8位，请重新输入"})
		} else {
			c.JSON(200, gin.H{"status": "用户名已存在，请重新注册"})
		}
	})

	//登录
	r.GET("/login", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")
		var pass PASS
		pass.Password, pass.Username = password, username
		if log(username, password) {
			c.JSON(200, pass)
		} else {
			c.JSON(200, gin.H{"status": "用户名或密码错误，请重新输入"})
		}
	})

	//注销
	r.DELETE("/logout", func(c *gin.Context) {
		var pass PASS
		pass.Username = "nil"
		c.JSON(200, pass)
	})

	r.Run(":8080")
}
