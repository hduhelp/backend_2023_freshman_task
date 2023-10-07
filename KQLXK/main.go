package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TODO struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Done     bool   `json"done"`
	indexID  int    `gorm:"primary_key"`
}

type UserInfo struct {
	Username string
	Password string
}

type wetherlog struct {
	Username string
	wether   bool
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
		var logok wetherlog
		c.BindJSON(&todo)
		db, _ = getDBConnection("logok")
		db.Where("username = ?", todo.Username).First(&logok)
		defer db.Close()
		if logok.wether == true {
			db, _ = getDBConnection("todolist")
			db.Create(&todo)
			c.JSON(200, gin.H{"msg": "添加成功"})
			fmt.Printf("%v", todos)
			defer db.Close()
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
		if logok.wether == true {
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
		var newtodo TODO
		index, _ := strconv.Atoi(c.Param("index"))
		c.BindJSON(&todo)
		if wetherlogin(todo) == true {
			db, _ = getDBConnection("todolist")
			db.Where("indexID =?", index).First(&newtodo)
			newtodo.indexID = todo.indexID
			newtodo.Content = todo.Content
			newtodo.Done = todo.Done
			newtodo.Username = todo.Username
			db.Save(&newtodo)
			c.JSON(200, gin.H{"msg": "修改成功"})

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
		todo.indexID = index
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
		c.BindJSON(&user)
		db, err := getDBConnection("yhmmm")
		if err != nil {
			fmt.Printf("无法连接到数据库")
			return
		}
		fmt.Printf("%v", user)
		db.Create(&user)
		defer db.Close()
		u := wetherlog{user.Username, false}
		db, _ = getDBConnection("logok")
		db.Create(&u)

		c.JSON(200, gin.H{"注册成功": "welcome"})

	})

	//登录
	r.POST("/login", func(c *gin.Context) {
		var user UserInfo
		c.BindJSON(&user)
		db, _ = getDBConnection("yhmmm")
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
			db.Model(&user).Update("wether", true)
			c.JSON(200, gin.H{"message": "登陆成功"})
		}
		defer db.Close()

	})

	//退出登录
	r.POST("/logout", func(c *gin.Context) {
		var username string
		var user wetherlog
		var todo TODO
		todo.Username = username
		c.BindJSON(&username)
		if wetherlogin(todo) == true {
			db.Where("username = ?", username).First(&user)
			db.Model(&user).Update("wether", false)
			c.JSON(200, gin.H{"msg": "退出成功"})
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
	if logok.wether == true {
		ok = true
		return ok
	} else {
		ok = false
		return ok
	}

}
