package main

import (
	"TodoList/config"
	"TodoList/utils"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

type Database struct {
	Username string
	Password string
	Host     string
	Port     int
	DBname   string
}

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
	config.InitConfig()
	ms := Database{
		Username: viper.GetString("database.mysql.username"),
		Password: viper.GetString("database.mysql.password"),
		Host:     viper.GetString("database.mysql.host"),
		Port:     viper.GetInt("database.mysql.port"),
		DBname:   viper.GetString("database.mysql.dbname"),
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", ms.Username, ms.Password, ms.Host, ms.Port, ms.DBname)
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
	index := 1
	for rows.Next() {
		var todo Todo
		//遍历表中所有行的信息
		rows.Scan(&todo.ID, &todo.Content, &todo.Done)
		i = todo.ID
		index++
	}
	rs, err := stmt.Exec(i+1, todo.Content, todo.Done)
	if err != nil {
		return
	}
	//	插入index
	_, err = rs.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	Index = index
	defer stmt.Close()
	return
}

// 获取所有todo
func getAll() ([]Todo, error) {
	db, err := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, content, done FROM todolist")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Content, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// 删除todo
func del(todos []Todo, index int) ([]Todo, error) {
	db, err := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM todolist WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todos[index-1].ID)
	if err != nil {
		return nil, err
	}

	todoss := append(todos[:index-1], todos[index:]...)
	return todoss, nil
}

// 修改todo
func update(todos []Todo, todo Todo, index int) ([]Todo, error) {
	db, err := connectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE todolist SET content=?, done=? WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.Content, todo.Done, todos[index-1].ID)
	if err != nil {
		return nil, err
	}

	todoss := todos
	todoss[index-1].Done = todo.Done
	todoss[index-1].Content = todo.Content

	return todoss, nil
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
			"token":   token,
		}
		c.JSON(http.StatusOK, result)
	})

	// 添加todo
	r.POST("/todo", utils.JWTHandler(), func(c *gin.Context) {
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
	r.DELETE("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// 修改todo
	r.PUT("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// 获取todo
	r.GET("/todo", utils.JWTHandler(), func(c *gin.Context) {
		todos, err := getAll()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, todos)
	})

	// 查询todo
	r.GET("/todo/:index", utils.JWTHandler(), func(c *gin.Context) {
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
