package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TODO struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Done     bool   `json"done"`
	IndexID  int    `gorm:"index"`
}

type UserInfo struct {
	Username string
	Password string
}

type wetherlog struct {
	Username string `gorm:"column:username"`
	Wether   string `gorm:"column:ok"`
}

var todos []TODO

func main() {
	r := gin.Default()
	//在数据库中添加数据表
	db, err := getDBConnection("yhmmm")
	if err != nil {
		fmt.Printf("无法连接到数据库")
		return
	}
	db.AutoMigrate(&UserInfo{})
	defer db.Close()

	db, err = getDBConnection("todolist")
	if err != nil {
		fmt.Printf("无法连接到数据库")
	}
	db.AutoMigrate(&TODO{})
	defer db.Close()

	db, err = getDBConnection("logok")
	if err != nil {
		fmt.Printf("无法连接到数据库")
	}
	db.AutoMigrate(&wetherlog{})
	defer db.Close()

	//添加TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		var todo1 TODO
		var logok wetherlog
		c.BindJSON(&todo)
		db, _ = getDBConnection("logok")
		db.Where("username = ?", todo.Username).First(&logok)
		defer db.Close()
		if logok.Wether == "true" {
			db, _ = getDBConnection("todolist")
			err = db.Where("indexID =?", todo.IndexID).First(&todo1).Error
			if err != nil {
				if gorm.IsRecordNotFoundError(err) {
					fmt.Println("没有找到符合条件的记录")
				} else {
					fmt.Println("发生了其他错误：", err)
				}
				c.JSON(200, gin.H{"msg": "indexID重复"})
			} else {
				db.Create(&todo)
				c.JSON(200, gin.H{"msg": "添加成功"})
				fmt.Printf("%v", todos)
				defer db.Close()
			}

		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}
		defer db.Close()
	})

	//删除TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		var todo TODO
		var logok wetherlog
		c.BindJSON(&todo)
		db, _ = getDBConnection("logok")
		db.Where("username = ?", todo.Username).First(&logok)
		defer db.Close()
		if logok.Wether == "true" {
			index, _ := strconv.Atoi(c.Param("index"))
			result := db.Where("indexID =?", index).First(&todo)
			if result.RecordNotFound() {
				// 没有找到记录
				c.JSON(200, gin.H{"msg": "没找到该todo"})
			} else if result.Error != nil {
				// 处理其他错误
				c.JSON(200, gin.H{"msg": "出现不知错误"})

			} else {
				// 找到了记录
				db.Delete(&todo)
				c.JSON(200, gin.H{"msg": "删除成功"})
				defer db.Close()
			}

		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}
		defer db.Close()
	})

	//修改TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		var todo TODO

		index, _ := strconv.Atoi(c.Param("index"))
		c.BindJSON(&todo)

		if wetherlogin(todo) == true {
			db, _ = getDBConnection("todolist")

			db.Model(&todo).Where("index_id =? AND username =?", index, todo.Username).Update("content", todo.Content, "done", todo.Done)
			if err != nil {

				if gorm.IsRecordNotFoundError(err) {
					fmt.Println("没有找到符合条件的记录")
				} else {
					fmt.Println("发生了其他错误：", err)
				}
			} else {
				c.JSON(200, gin.H{"msg": "修改成功"})
			}

		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}

		defer db.Close()
	})

	//列出TODO
	r.GET("/todo", func(c *gin.Context) {
		var todo TODO
		var todolist []TODO
		c.JSON(200, todos)
		if wetherlogin(todo) == true {
			db, _ = getDBConnection("todolist")
			db.Where("username =?", todo.Username).Find(&todolist)
			c.JSON(200, todolist)
		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}
		defer db.Close()
	})

	//查询TODO
	r.GET("/todo/:index", func(c *gin.Context) {
		var username string
		var todo TODO
		c.BindJSON(&username)
		index, _ := strconv.Atoi(c.Param("index"))
		todo.Username = username
		todo.IndexID = index
		if wetherlogin(todo) == true {

			db, _ = getDBConnection("todolist")
			db.Where("indexID =? AND username =?", index, username).First(&todo)
		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}
		defer db.Close()
	})

	//注册用户
	r.POST("/register", func(c *gin.Context) {
		var user UserInfo
		var user2 UserInfo
		c.BindJSON(&user)
		db, err := getDBConnection("yhmmm")
		if err != nil {
			fmt.Printf("无法连接到数据库")
			return
		}
		if strings.Contains(user.Username, " ") || user.Username == "" {
			c.JSON(200, gin.H{"msg": "用户名不能为空且不能包含空格"})
			return
		} else {
			_ = db.Where("username =?", user.Username).First(&user2)
			if user2.Username != user.Username {
				fmt.Printf("%v", user)
				db.Create(&user)
				defer db.Close()
				u := wetherlog{user.Username, "false"}
				db, _ = getDBConnection("logok")
				db.Create(&u)
				defer db.Close()
				c.JSON(200, gin.H{"注册成功": "welcome"})
			} else {
				c.JSON(200, gin.H{"msg": "该用户名已被注册，请更换用户名"})
			}

		}

	})

	//登录
	r.POST("/login", func(c *gin.Context) {
		var user UserInfo
		var iflog wetherlog
		c.BindJSON(&user)
		db, _ = getDBConnection("yhmmm")
		iflog.Username = user.Username
		iflog.Wether = "false"
		err = db.Where(&user).First(&user).Error
		fmt.Printf("%v", err)
		defer db.Close()
		if err != nil {

			if gorm.IsRecordNotFoundError(err) {
				fmt.Println("没有找到符合条件的记录")
			} else {
				fmt.Println("发生了其他错误：", err)
			}

			c.JSON(200, gin.H{"message": "用户名或密码错误"})
		} else {
			db, _ = getDBConnection("logok")
			db.Model(&iflog).Update("ok", "true")
			c.JSON(200, gin.H{"message": "登陆成功"})
		}
		defer db.Close()

	})

	//退出登录
	r.POST("/logout", func(c *gin.Context) {
		var username string
		var user wetherlog
		var todo TODO
		c.BindJSON(&username)
		todo.Username = username
		if wetherlogin(todo) == true {
			db, _ = getDBConnection("logok")
			err = db.Where("username = ?", username).First(&user).Error
			if err != nil {
				if gorm.IsRecordNotFoundError(err) {
					fmt.Println("没有找到符合条件的记录")
				} else {
					fmt.Println("发生了其他错误：", err)
				}
			} else {
				fmt.Printf("%v", user)
				db.Model(&user).Update("ok", "false")
				c.JSON(200, gin.H{"msg": "退出成功"})
			}

		} else {
			c.JSON(200, gin.H{"msg": "请先登录"})
		}
		defer db.Close()
	})

	r.Run(":8000")
}

// 根据数据库名获取数据库连接
func getDBConnection(database string) (*gorm.DB, error) {
	dsn := "root:yuchao@(localhost:3306)/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// 判断登陆状态
func wetherlogin(todo TODO) (ok bool) {
	var logok wetherlog
	db, _ := getDBConnection("logok")
	db.Where("username = ?", todo.Username).First(&logok)
	defer db.Close()
	if logok.Wether == "true" {
		ok = true
		return ok
	} else {
		ok = false
		return ok
	}

}
