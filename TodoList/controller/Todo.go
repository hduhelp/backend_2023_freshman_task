package controller

import (
	"log"
)

type Database struct {
	Username string
	Password string
	Host     string
	Port     int
	DBname   string
}

type Todo struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
	Done    string `json:"done"`
}

var todos []Todo

type User struct {
	ID       int    `json:"id"       form:"id"`
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// 添加todo
func Add(user string, todo Todo) (Index int, err error) {
	db, _ := ConnectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO todolist(id, usr, content, done) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	rows, err := db.Query("SELECT id, usr, content, done FROM todolist")
	if err != nil {
		log.Fatal(err.Error())
	}
	i := 1
	index := 1
	for rows.Next() {
		var todo Todo
		//遍历表中所有行的信息
		rows.Scan(&todo.ID, &todo.User, &todo.Content, &todo.Done)
		i = todo.ID
		index++
	}
	rs, err := stmt.Exec(i+1, user, todo.Content, todo.Done)
	if err != nil {
		return
	}
	//	插入index
	_, err = rs.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	Index = index
	defer stmt.Close()
	return
}

// 获取所有todo
func GetAll(user string) ([]Todo, error) {
	db, err := ConnectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, usr, content, done FROM todolist WHERE usr = ?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.User, &todo.Content, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// 删除todo
func Del(todos []Todo, index int) ([]Todo, error) {
	db, err := ConnectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM todolist WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todos[index-1].ID)
	if err != nil {
		return nil, err
	}

	todoss := append(todos[:index-1], todos[index:]...)
	return todoss, nil
}

// 修改todo
func Update(todos []Todo, todo Todo, index int) ([]Todo, error) {
	db, err := ConnectToDatabase() // 与 MySQL 数据库建立连接
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE todolist SET content=?, done=? WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.Content, todo.Done, todos[index-1].ID)
	if err != nil {
		return nil, err
	}

	todoss := todos
	todoss[index-1].Done = todo.Done
	todoss[index-1].Content = todo.Content

	return todoss, nil
}
