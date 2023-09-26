package models

import (
	"github.com/google/uuid"
	"regexp"
	"sync"
	"todo/util"
)

var Tokens sync.Map

func AddUser(username string, password string) (StatusCode, User) {
	var user User
	password = util.Str2md5(password + Env.GetString("salt"))
	res := db.Model(&User{}).Where("username = ?", username).Limit(1).Find(&user)
	if res.RowsAffected != 0 {
		Logger.Warning(res.Error)
		return UserExistsError, User{}
	} else {
		user = User{Username: username, Password: password, Role: false}
		if err := db.Model(&User{}).Create(&user).Error; err != nil {
			Logger.Error(err)
			return DatabaseError, User{}
		}
		return RegisterSuccess, user
	}
}

func CheckAuth(username string, password string) (StatusCode, User) {
	var user User
	password = util.Str2md5(password + Env.GetString("salt"))
	res := db.Model(&User{}).Where("username = ?", username).Find(&user).Limit(1)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return UserNotExistsError, User{}
	} else {
		res := db.Model(&User{}).Where(&User{Username: username, Password: password}).Find(&user).Limit(1)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return LoginFailed, User{}
		}
		return LoginSuccess, user
	}
}

func ChangePWD(username string, old string, new string) (StatusCode, User) {
	var user User
	old = util.Str2md5(old + Env.GetString("salt"))
	new = util.Str2md5(new + Env.GetString("salt"))
	res := db.Model(&User{}).Where(&User{Username: username, Password: old}).Limit(1).Find(&user)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return ChangePWDFailed, User{}
	} else {
		res := db.Model(&User{}).Where(&User{Username: username, Password: old}).Update("password", new)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ChangePWDFailed, User{}
		}
		res = db.Model(&User{}).Where(&User{Username: username, Password: new}).Find(&user).Limit(1)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ChangePWDFailed, User{}
		}
		return ChangePWDSuccess, user
	}
}

func SetEmail(username string, email string) (StatusCode, User) {
	var user User
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	if reg.MatchString(email) {
		res := db.Model(&User{}).Where("username = ?", username).Find(&user).Limit(1)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return UserNotExistsError, User{}
		} else {
			res := db.Model(&User{}).Where(&User{Username: username}).Update("email", email)
			if res.RowsAffected != 1 {
				Logger.Warning(res.Error)
				return SetEmailFailed, User{}
			}
			res = db.Model(&User{}).Where(&User{Username: username}).Find(&user).Limit(1)
			if res.RowsAffected != 1 {
				Logger.Warning(res.Error)
				return SetEmailFailed, User{}
			}
			return SetEmailSuccess, user
		}
	} else {
		Logger.Warning(EmailFormatError.Description)
		return EmailFormatError, User{}
	}
}

func GetUserByName(username string) (StatusCode, User) {
	var user User
	res := db.Model(&User{}).Where("username = ?", username).Limit(1).Find(&user)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return UserNotExistsError, User{}
	}
	return GetUserSuccess, user
}

func CreateToken(username string) (token string) {
	token = uuid.New().String()
	Tokens.Store(username, token)
	return
}
