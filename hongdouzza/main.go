package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
	"time"
)

// 日期time.TiME怎么处理，cookie ,session
var db *gorm.DB
var secretKey = []byte("your-secret-key")

type UserInfo struct {
	UserName1 string `json:"userName1"`
	Password  string `json:"password"`
}

// 结构体

type TODO struct {
	Number    int       `json:"number"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	Deadline  time.Time `json:"deadline"`
	Reminder  time.Time `json:"reminder"`
	Priority  int       `json:"priority"`
	UserName2 string    `json:"user_name2"`
}

// 连接数据库
func init() {

	var err error
	db, err = gorm.Open("mysql", "root:973100@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("无法连接数据库")
	}
	//绑定模型
	db.AutoMigrate(&TODO{}, &UserInfo{})
}

func createUser(username, password string) error {
	//创建一个新用户
	newUser := UserInfo{
		UserName1: username,
		Password:  password,
	}
	if err := db.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}

func main() {

	r := gin.Default()
	store := cookie.NewStore([]byte("session-secret"))
	r.Use(sessions.Sessions("mysession", store))
	//注册路由
	r.POST("/todolist_register", func(c *gin.Context) {
		// 启用 Session 中间件
		var inputUser UserInfo
		if err := c.BindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求"})
			return
		}

		// 检查用户名是否已存在
		var existingUser UserInfo
		if err := db.Where("user_name1=?", inputUser.UserName1).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}

		// 创建用户
		if err := createUser(inputUser.UserName1, inputUser.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建用户"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
		c.Redirect(http.StatusSeeOther, "/todolist_login")
	})

	//登录路由
	r.POST("/todolist_login", func(c *gin.Context) {
		var inputUser UserInfo
		if err := c.BindJSON(&inputUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求"})
			return
		}
		var dbUser UserInfo
		if err := db.Where("user_name1=?", inputUser.UserName1).First(&dbUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		if inputUser.Password != dbUser.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		// 登录成功，创建一个令牌
		token := jwt.New(jwt.SigningMethodHS256)

		// 添加一些声明 (claims)?
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 令牌过期时间（24小时）
		claims["username"] = inputUser.UserName1
		// 使用密钥签名令牌并获取字符串表示
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
			return
		}
		//将令牌存储在会话中
		session := sessions.Default(c)
		session.Set("token", tokenString)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": tokenString})
		//等待一会儿再重定向到todo清单
		time.Sleep(4 * time.Second)
		c.Redirect(http.StatusSeeOther, "/todo")
	})

	// 需要身份验证的路由
	r.GET("/todo-auth", AuthMiddleware(), func(c *gin.Context) {
		// 这里放置需要身份验证的代码，可以使用 c.Get("username") 来获取用户名
		username := c.GetString("username")
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated User: " + username})
	})
	// 增加 todo
	r.POST("/todo", AuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")
		var todo TODO
		_ = c.BindJSON(&todo)
		todo.UserName2 = username
		if err := db.Create(&todo).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "添加成功！"})

	})

	//删除 todo
	r.DELETE("/todo/:number", AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		number, _ := strconv.Atoi(c.Param("number"))
		var todo TODO
		if err := db.Where("user_name2=? AND number=?", username, number).First(&todo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "没有这个todo哦"})
			return
		}
		if err := db.Where("user_name2=? AND number=?", username, number).Delete(&todo).Error; err != nil {

			c.JSON(http.StatusNotFound, gin.H{"error": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "删除成功！"})

	})

	//更新 todo

	r.PUT("/todo/:number", AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		number, _ := strconv.Atoi(c.Param("number"))
		var todo TODO
		if err := db.Where("user_name2=? AND number=?", username, number).First(&todo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "无效的number"})
			return
		}

		var NewTodo TODO
		NewTodo.UserName2 = username
		NewTodo.Number = number
		if err := c.ShouldBindJSON(&NewTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Where("user_name2=? AND number=?", username, number).Model(&todo).Updates(NewTodo)
		c.JSON(200, gin.H{"status": "修改成功！"})

	})

	//列出 todo，有分页功能，并且按任务优先级为主要关键字，ddl为次要关键字排序

	r.GET("/todo", AuthMiddleware(), AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)

		var todos []TODO
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		offset := (page - 1) * pageSize

		if err := db.Where("user_name2=?", username).Order("priority DESC,deadline DESC").Limit(pageSize).Offset(offset).Find(&todos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)

	})

	//查询 todo
	r.GET("/todo/:number", AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		number, _ := strconv.Atoi(c.Param("number"))
		var todo TODO
		if err := db.Where("user_name2=? AND number=?", username, number).First(&todo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "没有这个todo哦"})
			return
		}
		c.JSON(200, todo)

	})

	r.Run()

}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从会话中获取令牌
		session := sessions.Default(c)
		tokenString := session.Get("token")
		if tokenString == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未经授权的访问"})
			c.Abort()
			return
		}

		// 解析 JWT 令牌
		token, err := jwt.Parse(tokenString.(string), func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			c.Set("username", username)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()

		}
		expTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
			c.Abort()
			return
		}

		// 令牌验证通过，可以继续执行需要身份验证的代码
		username := claims["username"].(string)

		// 这里放置需要身份验证的代码
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated User: " + username})

	}
}
