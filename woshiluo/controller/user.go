//
// user.go
// Copyright (C) 2023 Woshiluo Luo <woshiluo.luo@outlook.com>
//
// Distributed under terms of the GNU AGPLv3+ license.
//

package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"todolist/models"
)

func NewUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	if err := models.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

type UpdateUserData struct {
	Token string      `json:"token" binding:"required"`
	User  models.User `json:"user" binding:"required"`
}

func UpdateUser(c *gin.Context) {
	var origin_user models.User
	var id = c.Param("id")
	var data UpdateUserData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByToken(data.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Wrong Token"})
		return
	}

	if err := models.Db.First(&origin_user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "User Not Found"})
		return
	}

	if user.ID != origin_user.ID {
		c.JSON(http.StatusForbidden, "Forbidden")
		return
	}

	if err := models.Db.Model(&origin_user).Updates(data.User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, origin_user)
}

type NewTokenData struct {
	Auth Auth `json:"auth" binding:"required"`
}

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetUserByPassword(auth Auth) (models.User, error) {
	var user models.User

	if err := models.Db.First(&user, models.User{Username: auth.Username}).Error; err != nil {
		return models.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(user.Password)) != nil {
		return models.User{}, errors.New("Wrong username or password.")
	}

	return user, nil
}

func GetUserByToken(token_string string) (models.User, error) {
	var user models.User
	var token models.Token

	if err := models.Db.First(&token, models.Token{Token: token_string}).Error; err != nil {
		return models.User{}, err
	}

	if err := models.Db.First(&user, token.UserID).Error; err != nil {
		return models.User{}, errors.New("Wrong username or password.")
	}

	// Update `UpdatedTime`
	models.Db.Save(&token)

	return user, nil
}

func NewToken(c *gin.Context) {
	var data NewTokenData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := GetUserByPassword(data.Auth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	token := &models.Token{Token: uuid.New().String(), UserID: user.ID}

	if err := models.Db.Create(&token).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func DeleteToken(c *gin.Context) {
	var token models.Token
	var origin_token = c.Param("token")
	var data DeleteTodoData

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := models.Db.First(&token, models.Token{ Token: origin_token }).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Token not found"})
		return
	}

	if err := models.Db.Delete(&token).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}
