package main

import (
	"TodoList/utils"
	"database/sql"
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

const (
	username = "root"
	password = "xxxxxxxx"
	host     = "localhost"
	port     = 3316
	dbname   = "xxxxx"
)

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    string `json:"done"`
}

var todos []Todo

type User struct {
	ID       int    `json:"id"       form:"id"`
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// 连接数据库
func connectToDatabase() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, dbname)
	db, err := sql.Open("mysql", dataSourceName) // 打开mysql数据库
	if err != nil {
		return nil, err
	}
	err = db.Ping() // 检查连接是否建立，以确保连接存活
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 添加todo
func add(todo Todo) (Index int, err error) {
	db, _ := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO todolist(id, content, done) VALUES (?, ?, ?)")
	if err != nil {
		return
	}
	rows, err := db.Query("SELECT id, content, done FROM todolist")
	if err != nil {
		log.Fatal(err.Error())
	}
	i := 1
	for rows.Next() {
		var todo Todo
		//遍历表中所有行的信息
		rows.Scan(&todo.ID, &todo.Content, &todo.Done)
		i++
	}
	rs, err := stmt.Exec(i, todo.Content, todo.Done)
	if err != nil {
		return
	}
	//	插入index
	index, err := rs.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	Index = int(index)
	defer stmt.Close()
	return
}

// 获取所有todo
func getAll() (todos []Todo, err error) {
	db, _ := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, content, done FROM todolist")
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var todo Todo
		//遍历表中所有行的信息
		rows.Scan(&todo.ID, &todo.Content, &todo.Done)
		//将user添加到users中
		todos = append(todos, todo)
	}
	//最后关闭连接
	defer rows.Close()
	return
}

// 删除todo
func del(todos []Todo, index int) (todoss []Todo, err error) {
	db, _ := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM todolist WHERE id=?")
	if err != nil {
		log.Fatalln(err)
	}
	stmt.Exec(todos[index-1].ID)
	todoss = append(todos[:index-1], todos[index:]...)
	defer stmt.Close()
	return
}

// 修改todo
func update(todos []Todo, todo Todo, index int) (todoss []Todo, err error) {
	db, _ := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE todolist SET content=?, done=? WHERE id=?")
	if err != nil {
		log.Fatalln(err)
	}
	stmt.Exec(todo.Content, todo.Done, todos[index-1].ID)
	todoss = todos
	todoss[index-1].Done = todo.Done
	todoss[index-1].Content = todo.Content
	defer stmt.Close()
	return
}

// JWTHandler jwt拦截器
func JWTHandler() gin.HandlerFunc {
	return func(context *gin.Context) {

		//引入jwt实现登录后的会话记录,登录会话发生登录完成之后
		//header获取token
		session := sessions.Default(context)
		token := session.Get("username").(string)
		if token == "" {
			context.String(302, "请求未携带token无法访问!")
			context.Abort()
		}
		//解析token
		claims, err := utils.ParserUserToken(token)
		if claims == nil || err != nil {
			context.String(401, "未携带有效token或已过期")
			context.Abort()
		} else {
			//context.Set("user", claims.Username)
			context.Next()

		}
	}
}

func main() {
	db, err := connectToDatabase() // 与 MySQL 数据库建立连接
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
		row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", usrname)
		err = row.Scan(&user.ID, &user.UserName, &user.Password)
		if user.ID == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户不存在",
			})
			return
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
			"data":    user,
			"token":   token,
		}
		c.JSON(http.StatusOK, result)
	})

	// 添加todo
	r.POST("/todo", JWTHandler(), func(c *gin.Context) {
		var todo Todo
		c.BindJSON(&todo)
		index, err := add(todo)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"index":  index,
			"status": "success",
		})
	})

	// 删除todo
	r.DELETE("/todo/:index", JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		todos, _ := getAll()
		if index-1 >= len(todos) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "此位置的todo不存在",
			})
			return
		}
		todos, err = del(todos, index)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, todos)
	})

	// 修改todo
	r.PUT("/todo/:index", JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		var todo Todo
		c.BindJSON(&todo)
		todos, _ := getAll()
		if index-1 >= len(todos) {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "此位置的todo不存在",
			})
			return
		}
		todos, err = update(todos, todo, index)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, todos)
	})

	// 获取todo
	r.GET("/todo", JWTHandler(), func(c *gin.Context) {
		todos, err := getAll()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, todos)
	})

	// 查询todo
	r.GET("/todo/:index", JWTHandler(), func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		todos, err := getAll()
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
