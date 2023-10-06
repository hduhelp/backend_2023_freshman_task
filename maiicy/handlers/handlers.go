package handlers

import (
	"login-system/db_handle"
	"login-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context) {
	var newUser db_handle.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db_handle.GetUserByUsername(newUser.Username)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{"message": "该用户名已被注册"})
		return
	}

	if !utils.IsValidUsername(newUser.Username) {
		c.JSON(http.StatusCreated, gin.H{"message": "该用户名不合法"})
		return
	}

	md5PassWord, _ := utils.CalculateMD5(newUser.Password)
	err = db_handle.InsertUser(newUser.Username, md5PassWord, newUser.Info)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

func LoginHandler(c *gin.Context) {
	var loginData loginData

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db_handle.GetUserByUsername(loginData.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户名不存在"})
		return
	}

	md5Password, err := utils.CalculateMD5(loginData.Password)
	if err != nil {
		return
	}

	if md5Password != user.Password {
		c.JSON(http.StatusOK, gin.H{"message": "密码不正确"})
		return
	}

	userData := utils.User{ID: user.ID, Username: user.Username, Password: user.Password}

	token, err := utils.GenerateJWT(userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	c.Header("Authorization", token)

}

func LogoutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	err := db_handle.InsertJWTIntoBlacklist(tokenString)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户已成功注销",
	})
}

type addTodoData struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

func TodoAddHandler(c *gin.Context) {
	var newTodo addTodoData

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)

	err = db_handle.InsertTodo(user.ID, newTodo.Title, newTodo.Date)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "用户已成功添加任务",
	})
}

type delTodoData struct {
	TodoID int `json:"todo_id"`
}

func TodoDelHandler(c *gin.Context) {
	var delTodo delTodoData

	if err := c.ShouldBindJSON(&delTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return
	}

	todo, err := db_handle.FindTodoByID(delTodo.TodoID)
	if err != nil {
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户不存在此任务",
		})
		return
	}

	err = db_handle.DeleteTodo(delTodo.TodoID)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "用户已成功删除任务",
	})
}

type updateTodoData struct {
	TodoID    int    `json:"todo_id"`
	Title     string `json:"title"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}

func TodoUpdateHandler(c *gin.Context) {
	var updateTodo updateTodoData

	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return
	}

	todo, err := db_handle.FindTodoByID(updateTodo.TodoID)
	if err != nil {
		return
	}
	if todo.UserID != user.ID {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户不存在此任务",
		})
		return
	}

	err = db_handle.UpdateTodo(updateTodo.TodoID, updateTodo.Title, updateTodo.Completed, updateTodo.Date)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "用户已成功更新任务",
	})
}

func GetAllTodoHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return
	}

	todos, err := db_handle.FindTodosByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取Todo项"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": todos})
}

func GetIDTodoHandler(c *gin.Context) {
	todoParam := c.Param("id")
	todoID, err := strconv.Atoi(todoParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式不正确"})
		return
	}
	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return
	}
	todo, err := db_handle.FindTodoByID(todoID)
	if todo.UserID != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不存在对应的任务"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "获取成功", "data": todo})

}

func GetDateTodoHandler(c *gin.Context) {
	date := c.Param("date")
	dateData, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式不正确"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		return
	}

	todos, err := db_handle.FindTodosBeforeTime(dateData, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取Todo项"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": todos})
}
