package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	Id       int
	Username string
	Password string
}

func initDB() (err error) {
	dsn := "root:My_Shit_SQL123@(127.0.0.1)/go_db?charset=utf8mb4&parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func user_insert(Id int, Username string, Password string) {
	s := "insert into user_tbl(Id, Username, Password) values(?,?,?)"
	r, err := db.Exec(s, Id, Username, Password)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}
}

func user_insertManyRow() {
	for _, user := range users {
		user_insert(user.Id, user.Username, user.Password)
	}
}

func user_queryOneRow() {
	s := "select * from user_tbl where id = ?"
	var user User
	err := db.QueryRow(s, 1).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		fmt.Println("u: %v\n", user)
	}
}

/*
func user_queryManyRow() {
	s := "select * from user_tbl "
	r, err := db.Query(s)
	var user User
	defer r.Close()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		for r.Next() {
			r.Scan(&user.Id, &user.Username, &user.Password)
			fmt.Printf("u:%v\n", user)
		}
	}

}
*/

func user_queryManyRow() []User {
	s := "select Id, Username, Password from user_tbl"
	rows, err := db.Query(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		users = append(users, user)
	}
	return users
}

func user_update(username string, password string, id int) {
	s := "update user_tbl set username=?, password=? where id=?"
	r, err := db.Exec(s, username, password, id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.RowsAffected()
		fmt.Printf("i: %v\n", i)
	}
}

func user_delete() {
	s := "delete from user_tbl where id=?"
	r, err := db.Exec(s, 2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}

}

func user_deleteAll() {
	s := "delete from user_tbl"
	_, err := db.Exec(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
