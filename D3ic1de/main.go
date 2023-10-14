package main

import (
	"TodoList/controller"
	"TodoList/utils"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
	Done    string `json:"done"`
}

var todos []Todo

type User struct {
	ID       int    `json:"id"       form:"id"`
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func main() {
	db, err := controller.ConnectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 延迟执行，确保在 main 函数退出时关闭数据库连接
	r := gin.Default()
	// 创建cookie存储
	store := cookie.NewStore([]byte("secret"))
	//路由上加入session中间件
	r.Use(sessions.Sessions("token", store))
	// 注册
	r.POST("/register", func(c *gin.Context) {
		//获取参数
		usrname := c.PostForm("username")
		passwd := c.PostForm("password")

		//数据验证
		if len(usrname) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		var user User
		row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", usrname)
		err = row.Scan(&user.ID, &user.UserName, &user.Password)
		if user.ID != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户已存在",
			})
			return
		}
		//创建用户
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    500,
				"message": "密码加密错误",
			})
			return
		}
		newUser := User{
			UserName: usrname,
			Password: string(hasedPassword),
		}
		sql := "INSERT INTO users(username, password) VALUES(?, ?)"
		db.Exec(sql, newUser.UserName, newUser.Password)

		//返回结果
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注册成功",
		})
	})

	// 登录
	r.POST("/login", func(c *gin.Context) {
		var (
			user   User  // 用户结构体，用于存储结果查询的记录
			result gin.H // Gin框架使用的Map集合类型，用于将结果渲染为 JSON 格式并发送给客户端
		)
		usrname := c.PostForm("username")
		passwd := c.PostForm("password")
		// 进行用户验证
		// 执行 SQL 查询，并返回*sql.Row对象，其中包含结果集的单行记录
		if len(usrname) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", usrname)
		err = row.Scan(&user.ID, &user.UserName, &user.Password)
		if err != nil {
			if user.ID == 0 {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"code":    422,
					"message": "用户不存在",
				})
				return
			}
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码错误",
			})
			return
		}
		token := utils.MakeUserToken(usrname)
		session := sessions.Default(c)
		session.Set("username", token)
		session.Save()
		//返回结果
		result = gin.H{
			"code":    200,
			"message": "登陆成功",
			"token":   token,
		}
		c.JSON(http.StatusOK, result)
	})

	// 添加todo
	r.POST("/todo", utils.JWTHandler(), func(c *gin.Context) {
		user := c.GetString("user")
		var todo Todo
		c.BindJSON(&todo)
		index, err := controller.Add(user, controller.Todo(todo))
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"code":   200,
			"index":  index,
			"status": "success",
		})
	})

	// 删除todo
	r.DELETE("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		user := c.GetString("user")
		todos, _ := controller.GetAll(user)
		if index-1 >= len(todos) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "此位置的todo不存在",
			})
			return
		}
		todos, err = controller.Del(todos, index)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"message": "success",
		})
	})

	// 修改todo
	r.PUT("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		user := c.GetString("user")
		var todo Todo
		c.BindJSON(&todo)
		todos, _ := controller.GetAll(user)
		if index-1 >= len(todos) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "此位置的todo不存在",
			})
			return
		}
		todos, err = controller.Update(todos, controller.Todo(todo), index)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
		})
	})

	// 获取todo
	r.GET("/todo", utils.JWTHandler(), func(c *gin.Context) {
		user := c.GetString("user")
		todos, err := controller.GetAll(user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Sprintf(user)
		c.JSON(200, todos)
	})

	// 查询todo
	r.GET("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		user := c.GetString("user")
		todos, err := controller.GetAll(user)
		if err != nil {
			log.Fatal(err)
		}
		if index-1 >= len(todos) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "此位置的todo不存在",
			})
			return
		}
		c.JSON(200, todos[index-1])
	})

	r.Run(":8080")
}
