package model

import "fmt"

func migration() {
	// 自动迁移模式
	err := DB.AutoMigrate(&User{})
	err2 := DB.AutoMigrate(&Task{})
	if err != nil {
		fmt.Println("用户数据表迁移失败")
	}
	if err2 != nil {
		fmt.Println("todolist数据表迁移失败")
	}
}
