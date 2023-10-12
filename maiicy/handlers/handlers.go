package handlers

import (
	"login-system/db_handle"
	"login-system/requests"
	"login-system/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHandler(c *gin.Context) {
	var RegisterData requests.RegisterRequest

	if err := c.ShouldBindJSON(&RegisterData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误", "error": err.Error()})
		return
	}

	if !utils.IsValidUsername(RegisterData.Username) {
		c.JSON(http.StatusCreated, gin.H{"message": "用户名不合法"})
		return
	}

	// 检查用户名是否已经存在
	_, err := db_handle.GetUserByUsername(RegisterData.Username)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{"message": "用户名已被注册"})
		return
	}

	md5Password, err := utils.CalculateMD5(RegisterData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	if err := db_handle.InsertUser(RegisterData.Username, md5Password, RegisterData.Info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "用户注册失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

func LoginHandler(c *gin.Context) {
	var loginData requests.LoginRequest

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误", "error": err.Error()})
		return
	}

	user, err := db_handle.GetUserByUsername(loginData.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户名不存在"})
		return
	}

	md5Password, err := utils.CalculateMD5(loginData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	if md5Password != user.Password {
		c.JSON(http.StatusOK, gin.H{"message": "密码不正确"})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func LogoutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	err := db_handle.InsertJWTIntoBlacklist(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注销失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户已成功注销",
	})
}

func TodoAddHandler(c *gin.Context) {
	var newTodoData requests.AddTodoRequest

	if err := c.ShouldBindJSON(&newTodoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误", "error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	err = db_handle.InsertTodo(user.ID, newTodoData.Title, newTodoData.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已成功添加"})
}

func TodoDelHandler(c *gin.Context) {
	todoParam := c.Param("id")
	todoID, err := strconv.Atoi(todoParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID格式不正确"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	todo, err := db_handle.FindTodoByID(todoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	// 验证用户是否有权限删除该待办事项
	if todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "没有权限删除该任务"})
		return
	}

	// 删除待办事项
	err = db_handle.DeleteTodo(todoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已成功删除"})
}

func TodoUpdateHandler(c *gin.Context) {
	var updateTodoData requests.UpdateTodoRequest

	if err := c.ShouldBindJSON(&updateTodoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误", "error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	todo, err := db_handle.FindTodoByID(updateTodoData.TodoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "没有权限更新该任务"})
		return
	}

	err = db_handle.UpdateTodo(updateTodoData.TodoID, updateTodoData.Title, updateTodoData.Completed, updateTodoData.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务已成功更新"})
}

func GetAllTodoHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	todos, err := db_handle.FindTodosByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取待办事项"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": todos})
}

func GetIDTodoHandler(c *gin.Context) {
	todoParam := c.Param("id")
	todoID, err := strconv.Atoi(todoParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID格式不正确"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	todo, err := db_handle.FindTodoByID(todoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "内部错误"})
		return
	}

	if todo.UserID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"message": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": todo})
}

func GetDateTodoHandler(c *gin.Context) {
	date := c.Param("date")
	dateData, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "日期格式不正确"})
		return
	}

	tokenString := c.GetHeader("Authorization")
	user, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	todos, err := db_handle.FindTodosBeforeTime(dateData, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取待办事项"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": todos})
}
