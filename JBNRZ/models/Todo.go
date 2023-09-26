package models

import (
	"encoding/base64"
	"time"
)

func AddTodo(name string, detail string, endTime int64, username string) (StatusCode, Todo) {
	var todo Todo
	detail = base64.StdEncoding.EncodeToString([]byte(detail))
	res := db.Model(&Todo{}).Where("item_name = ?", name).Limit(1).Find(&todo)
	endTime = time.Now().Unix() + endTime
	if res.RowsAffected != 0 {
		Logger.Warning(res.Error)
		return ItemExistsError, Todo{}
	} else {
		todo = Todo{Username: username, ItemName: name, Detail: detail, EndTime: endTime}
		if err := db.Model(&Todo{}).Create(&todo).Error; err != nil {
			Logger.Error(err)
			return DatabaseError, Todo{}
		}
		return AddTodoSuccess, todo
	}
}

func DelTodo(name string, username string) (StatusCode, Todo) {
	var todo Todo
	res := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).Find(&todo).Limit(1)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return ItemNotExistsError, Todo{}
	} else {
		if err := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).
			Delete(&Todo{}).Error; err != nil {
			return DatabaseError, Todo{}
		}
		if err := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).
			Limit(1).Find(&Todo{}).Error; err != nil {
			Logger.Error(err)
			return DatabaseError, Todo{}
		}
		return DelTodoSuccess, Todo{}
	}
}

func GetTodo(name string, username string) (StatusCode, Todo) {
	var todo = Todo{}
	res := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).Find(&todo).Limit(1)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return ItemNotExistsError, Todo{}
	} else {
		todo = Todo{ItemName: name, Username: username}
		res = db.Model(&Todo{}).Where(&Todo{ItemName: name, Username: username}).Find(&todo).Limit(1)
		if res.Error != nil {
			Logger.Error(res.Error)
			return DatabaseError, Todo{}
		}
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ItemNotExistsError, Todo{}
		}
		return ListTodoSuccess, todo
	}
}

func ChangeTodo(name string, username string, endTime int64) (StatusCode, Todo) {
	var todo Todo
	res := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).Limit(1).Find(&todo)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return ItemNotExistsError, Todo{}
	} else {
		todo = Todo{ItemName: name, Username: username}
		res = db.Model(&Todo{}).Where(&todo).Find(&todo).Update("end_time", time.Now().Unix()+endTime)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ChangeTodoFailed, Todo{}
		}
		res = db.Model(&Todo{}).Where(&Todo{ItemName: name, Username: username}).Limit(1).Find(&todo)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ListTodoFailed, Todo{}
		}
		return ChangeTodoSuccess, todo
	}
}

func ListTodo(username string, from int, limit int) (StatusCode, []Todo) {
	var todos []Todo
	res := db.Model(&Todo{}).Where("username = ?", username).
		Offset(from).Limit(limit).Find(&todos)
	if res.Error != nil {
		Logger.Warning(res.Error)
		return ListTodoFailed, []Todo{}
	}
	for k, todo := range todos {
		decoded, _ := base64.StdEncoding.DecodeString(todo.Detail)
		todos[k].Detail = string(decoded)
	}
	return ListTodoSuccess, todos
}

func ListAll(from int, limit int) (StatusCode, []Todo) {
	var todos []Todo
	res := db.Model(&Todo{}).Offset(from).Limit(limit).Find(&todos)
	if res.Error != nil {
		Logger.Warning(res.Error)
		return ListTodoFailed, []Todo{}
	}
	for k, todo := range todos {
		decoded, _ := base64.StdEncoding.DecodeString(todo.Detail)
		todos[k].Detail = string(decoded)
	}
	return ListTodoSuccess, todos
}

func Done(name string, username string) (StatusCode, Todo) {
	var todo Todo
	res := db.Model(&Todo{}).Where("item_name = ? AND username = ?", name, username).Limit(1).Find(&todo)
	if res.RowsAffected != 1 {
		Logger.Warning(res.Error)
		return ItemNotExistsError, Todo{}
	} else {
		todo = Todo{ItemName: name, Username: username}
		res = db.Model(&Todo{}).Where(&todo).Find(&todo).Update("finished", true)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ChangeTodoFailed, Todo{}
		}
		res = db.Model(&Todo{}).Where(&Todo{ItemName: name, Username: username}).Limit(1).Find(&todo)
		if res.RowsAffected != 1 {
			Logger.Warning(res.Error)
			return ListTodoFailed, Todo{}
		}
		return ChangeTodoSuccess, todo
	}
}
