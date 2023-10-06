package controller

import (
	"TODOlist/dao/mysql"
	"TODOlist/middlewares/jwt"
	"TODOlist/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var user models.User
	_ = c.BindJSON(&user)
	result := mysql.DB.Where("username = ? AND password = ?", user.Username, user.Password).Find(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "user does not exist",
		})
		return
	}
	token, _ := jwt.GenerateToken(user.UserID, user.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"user":    user,
		"token":   token,
	})
}

func Register(c *gin.Context) {
	var user models.User
	_ = c.BindJSON(&user)
	result := mysql.DB.Create(user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "register failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"user":    user,
	})
}

func ParseToken(c *gin.Context) {
	token := c.Query("token")
	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "login expired",
		})
	}
	fmt.Println(claims)
	c.JSON(http.StatusOK, gin.H{
		"userid":   claims.UserID,
		"username": claims.Username,
	})
}
