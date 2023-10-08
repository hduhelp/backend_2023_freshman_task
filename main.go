package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type TODO struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var todos []TODO

func main() {
	r := gin.Default()

	// 添加 TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		c.BindJSON(&todo)
		todos = append(todos, todo)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 删除 TODO
	r.DELETE("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		todos = append(todos[:index], todos[index+1:]...)
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 修改 TODO
	r.PUT("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		todos[index] = todo
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 获取 TODO
	r.GET("/todo", func(c *gin.Context) {
		c.JSON(200, todos)
	})

	// 查询 TODO
	r.GET("/todo/:index", func(c *gin.Context) {
		index, _ := strconv.Atoi(c.Param("index"))
		c.JSON(200, todos[index])
	})

	// 登录接口
	r.POST("/login", loginHandler)

	// 鉴权中间件
	r.Use(authMiddleware)

	// 获取个人信息接口
	r.GET("/user", getUserHandler)

	// 登出接口
	r.GET("/logout", logoutHandler)

	r.Run(":8080")

}

// 登录处理函数
func loginHandler(c *gin.Context) {
	// 解析请求参数
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: 校验用户名和密码是否正确
	// ...

	// 登录成功，生成并返回token
	token := "your_generated_token"
	c.JSON(200, gin.H{"token": token})
}

// 鉴权中间件
func authMiddleware(c *gin.Context) {
	// 解析token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// TODO: 校验token是否有效
	// ...

	// 鉴权通过，继续处理请求
	c.Next()
}

// 获取个人信息处理函数
func getUserHandler(c *gin.Context) {
	// TODO: 根据鉴权信息获取个人信息
	userInfo := gin.H{
		"name":  "John Doe",
		"email": "john.doe@example.com",
		// 其他个人信息...
	}

	c.JSON(200, userInfo)
}

// 登出处理函数
func logoutHandler(c *gin.Context) {
	// TODO: 注销token
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}
