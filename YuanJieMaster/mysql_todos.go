package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

/*
	type TODO struct {
		Id    int
		Content  string
		Done     bool
		Priority int
	}
*/
/*
func todo_insert(Id int, Content string, Done bool, Priority int) {
	s := "insert into todo_tbl(Id, Content, Done, Priority) values(?,?,?,?)"
	r, err := db.Exec(s, Id, Content, Done, Priority)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}
}
*/

func todo_insert(Id int, Content string, done int, Priority int) {
	s := "insert into todo_tbl(Id, Content, Done, Priority) values(?, ?, ?, ?)"
	r, err := db.Exec(s, Id, Content, done, Priority)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}
}

func todo_insertManyRow() {
	for _, todo := range todos {
		done := 0
		if todo.Done {
			done = 1
		}
		todo_insert(todo.Id, todo.Content, done, todo.Priority)
	}
}

func todo_queryOneRow() {
	s := "select * from todo_tbl where id = ?"
	var todo TODO
	err := db.QueryRow(s, 1).Scan(&todo.Id, &todo.Content, &todo.Done, &todo.Priority)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		fmt.Println("todo: %v\n", todo)
	}
}

/*
	func todo_queryManyRow() {
		s := "select * from todo_tbl "
		r, err := db.Query(s)
		var todo TODO
		defer r.Close()
		if err != nil {
			fmt.Printf("err:%v\n", err)
		} else {
			for r.Next() {
				r.Scan(&todo.Id, &todo.Content, &todo.Done, &todo.Priority)
				fmt.Printf("todo:%v\n", todo)
			}
		}
	}
*/
func todo_queryManyRow() []TODO {
	s := "select Id, Content, Done, Priority from todo_tbl"
	rows, err := db.Query(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	defer rows.Close()

	var todos []TODO
	for rows.Next() {
		var todo TODO
		err := rows.Scan(&todo.Id, &todo.Content, &todo.Done, &todo.Priority)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		todos = append(todos, todo)
	}
	return todos
}

func todo_update(id int, content string, done bool, priority int) {
	s := "update todo_tbl set Content=?, Done=?, Priority=?, where Index=?"
	r, err := db.Exec(s, content, done, priority, id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.RowsAffected()
		fmt.Printf("i: %v\n", i)
	}
}

func todo_delete() {
	s := "delete from todo_tbl where id=?"
	r, err := db.Exec(s, 2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}
}

func todo_deleteAll() {
	s := "delete from todo_tbl"
	_, err := db.Exec(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

/*
func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		fmt.Println("连接成功")
	}
}
*/
