package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type TODO struct {
	Id       int
	Content  string
	Done     bool
	Priority int
}

var todos []TODO
var everyday_todos []TODO
var users []User
var login_flag int = 0

func check_id(id int, c *gin.Context, todos []TODO) {
	if id == 0 {
		c.JSON(400, gin.H{"status": "error", "massage": "0 or not digit , please input digit in range"})
		return
	}
	if id < 0 || id > len(todos) {
		c.JSON(400, gin.H{"status": "error", "massage": "out of range , please input digit in range"})
		return
	}
}

func sort(sort_todos []TODO, n int) []TODO {
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if sort_todos[j].Priority > sort_todos[j+1].Priority {
				sort_todos[j], sort_todos[j+1] = sort_todos[j+1], sort_todos[j]
			}
		}
	}
	return sort_todos
}

func main() {

	r := gin.Default()

	err := initDB()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		fmt.Println("连接成功")
	}

	everyday_todos = everyday_todo_queryManyRow()
	todos = todo_queryManyRow()
	users = user_queryManyRow()

	//注册账号
	r.POST("/register", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		for i := 0; i < len(users); i++ {
			if users[i].Username == user.Username {
				c.JSON(400, gin.H{"status": "error", "massage": "the username has been registered"})
				return
			}
		}
		users = append(users, user)
		n := len(users)
		user_insert(n, user.Username, user.Password)
		c.JSON(200, gin.H{"status": "ok", "massage": "new account has registered"})
	})

	//注销账号
	r.DELETE("/delete_user", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		flag := 0
		for i := 0; i < len(users); i++ {
			if users[i].Username == user.Username {
				users = append(users[:i], users[i+1:]...)
				flag = 1
				break
			}
		}
		if flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "the username didn't exist"})
			return
		}
		for j := 0; j < len(users); j++ {
			users[j].Id = j + 1
		}
		user_deleteAll()
		user_insertManyRow()
		c.JSON(200, gin.H{"status": "ok", "massage": "the account has deleted"})
	})

	//登陆账号
	r.PUT("/login", func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		flag1 := 0
		target := -1
		for i := 0; i < len(users); i++ {
			if users[i].Username == user.Username {
				flag1 = 1
				target = i
				break
			}
		}
		if flag1 == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "the username didn't exist"})
			return
		}
		flag2 := 0
		for i := 0; i < len(users); i++ {
			if users[target].Password == user.Password {
				flag2 = 1
				break
			}
		}
		if flag2 == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "the password is wrong"})
			return
		}
		login_flag = 1
		c.JSON(200, gin.H{"status": "ok", "massage": "login successfully"})
	})

	//添加 TODO
	r.POST("/todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		var todo TODO
		c.BindJSON(&todo)
		todos = append(todos, todo)
		n := len(todos)
		todos = sort(todos, n)
		for i := 0; i < n; i++ {
			todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//加入每日 TODO
	r.POST("/add_everyday_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		for i := 0; i < len(everyday_todos); i++ {
			todos = append(todos, everyday_todos[i])
		}
		n := len(todos)
		todos = sort(todos, n)
		for i := 0; i < n; i++ {
			todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//删除 TODO
	r.DELETE("/delete_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		todos = append(todos[:id-1], todos[id:]...)
		n := len(todos)
		todos = sort(todos, n)
		for i := 0; i < n; i++ {
			todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//删除已完成 TODO
	r.DELETE("/delete_finished_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		for i := 0; i < len(todos); i++ {
			if todos[i].Done {
				todos = append(todos[:i], todos[i+1:]...)
				i--
			}
		}
		todos = sort(todos, len(todos))
		for j := 0; j < len(todos); j++ {
			todos[j].Id = j + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//删除全部 TODO
	r.DELETE("/delete_all_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		todos = todos[0:0]
		c.JSON(200, gin.H{"status": "ok"})
	})

	//修改 TODO
	r.PUT("/modify_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		var todo TODO
		c.BindJSON(&todo)
		todos[id-1] = todo
		n := len(todos)
		todos = sort(todos, n)
		for i := 0; i < n; i++ {
			todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//完成 TODO
	r.PUT("/finish_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		todos[id-1].Done = true
		n := len(todos)
		for i := 0; i < n; i++ {
			todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//列出 TODO
	r.GET("/todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		c.JSON(200, todos)
	})

	//查询 TODO
	r.GET("/todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		c.JSON(200, todos[id-1])
	})

	//保存 TODO
	r.PUT("/todo/save", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		todo_deleteAll()
		todo_insertManyRow()
		c.JSON(200, gin.H{"status": "ok", "massage": "save successfully"})
	})

	//添加每日 TODO
	r.POST("/everyday_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		var everyday_todo TODO
		c.BindJSON(&everyday_todo)
		everyday_todos = append(everyday_todos, everyday_todo)
		n := len(everyday_todos)
		everyday_todos = sort(everyday_todos, n)
		for i := 0; i < n; i++ {
			everyday_todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//删除每日 TODO
	r.DELETE("/delete_everyday_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		everyday_todos = append(everyday_todos[:id-1], everyday_todos[id:]...)
		n := len(everyday_todos)
		everyday_todos = sort(everyday_todos, n)
		for i := 0; i < n; i++ {
			everyday_todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//删除全部每日 TODO
	r.DELETE("/delete_all_everyday_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		everyday_todos = everyday_todos[0:0]
		c.JSON(200, gin.H{"status": "ok"})
	})

	//修改每日 TODO
	r.PUT("/modify_everyday_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		var everyday_todo TODO
		c.BindJSON(&everyday_todo)
		everyday_todos[id-1] = everyday_todo
		n := len(everyday_todos)
		everyday_todos = sort(everyday_todos, n)
		for i := 0; i < n; i++ {
			everyday_todos[i].Id = i + 1
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	//列出每日 TODO
	r.GET("/everyday_todo", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		c.JSON(200, everyday_todos)
	})

	//查询每日 TODO
	r.GET("/everyday_todo/:id", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		check_id(id, c, todos)
		c.JSON(200, everyday_todos[id-1])
	})

	//保存每日 TODO
	r.PUT("/everyday_todo/save", func(c *gin.Context) {
		if login_flag == 0 {
			c.JSON(400, gin.H{"status": "error", "massage": "please login first"})
			return
		}
		everyday_todo_deleteAll()
		everyday_todo_insertManyRow()
		c.JSON(200, gin.H{"status": "ok", "massage": "save successfully"})
	})

	r.Run(":8080")
}

//go env -w GOPROXY=https://goproxy.cn,direct
//go get -u github.com/gin-gonic/gin
