package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	JSON_SUCCESS int = 1
	JSON_ERROR   int = 0
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type (
	// 数据库模型
	todoModel struct {
		gorm.Model
		Title string `json:"title"`
		//apifox上没有bool类型，采用int接收
		Completed int `json:"completed"`
	}
	// 格式化输出
	fmtTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

// 命名数据表
func (todoModel) TableName() string {
	return "todos"
}

var db *gorm.DB

// 初始化
func init() {
	var err error
	//后续可以写到配置文件中
	dsn := "root:123456@tcp(127.0.0.1:3306)/whxxxxxxxxxx?charset=utf8mb4&parseTime=True&loc=Local"
	//链接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败，请检查数据库配置")
	} else {
		fmt.Println("数据库连接成功")
	}
	//根据结构体格式创建数据表
	db.AutoMigrate(&todoModel{})
}

func add(c *gin.Context) {

	//调试
	/* if c == nil {
		fmt.Println("c is nil")
	} */

	/* if c.PostForm("title") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "标题不能为空",
		})
	} */
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	//储存数据
	//result := db.Save(&todo)
	db.Save(&todo)
	/* if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
	} */
	c.JSON(http.StatusOK, gin.H{
		"status":     JSON_SUCCESS,
		"message":    "创建成功",
		"resourceId": todo.ID,
	})
}

// 获取所有条目
func all(c *gin.Context) {
	//切片获取
	var todos []todoModel
	//切片输出
	var _todos []fmtTodo
	db.Find(&todos)

	// 没有数据
	if len(todos) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "没有数据",
		})
		return
	}

	// 格式化
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		//添加到切片
		_todos = append(_todos, fmtTodo{
			ID:        item.ID,
			Title:     item.Title,
			Completed: completed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todos,
	})

}

// 根据id获取一个条目
func take(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  JSON_ERROR,
			"message": "对应任务不存在",
		})
		return
	}
	//完成状态
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}
	//返回数据

	_todo := fmtTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todo,
	})
}

// 更新一个条目
func update(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  JSON_ERROR,
			"message": "对应任务不存在",
		})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "人物更新成功",
	})
}

// 删除条目
func del(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "对应任务不存在",
		})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "删除成功!",
	})
}

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1/todo")
	{
		v1.POST("/", add)      // 添加新条目
		v1.GET("/", all)       // 查询所有条目
		v1.GET("/:id", take)   // 获取单个条目
		v1.PUT("/:id", update) // 更新单个条目
		v1.DELETE("/:id", del) // 删除单个条目
	}
	r.Run(":80")
}
