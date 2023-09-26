package models

import "todo/util"

func InitAdmin() {
	var user User
	username := "admin"
	password := util.Str2md5("admin" + Env.GetString("salt"))
	res := db.Model(&User{}).Where("username = ?", username).Limit(1).Find(&user)
	if res.RowsAffected != 0 {
		Logger.Warning(res.Error)
		return
	} else {
		user = User{Username: username, Password: password, Role: true}
		if err := db.Model(&User{}).Create(&user).Error; err != nil {
			Logger.Error(err)
			return
		}
		return
	}
}
