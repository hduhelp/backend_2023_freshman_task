package main

import "fmt"

func everyday_todo_insert(Id int, Content string, done int, Priority int) {
	s := "insert into everyday_todo_tbl(Id, Content, Done, Priority) values(?, ?, ?, ?)"
	r, err := db.Exec(s, Id, Content, done, Priority)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		i, _ := r.LastInsertId()
		fmt.Printf("i: %v\n", i)
	}
}

func everyday_todo_insertManyRow() {
	for _, everyday_todo := range everyday_todos {
		todo_insert(everyday_todo.Id, everyday_todo.Content, 0, everyday_todo.Priority)
	}
}

func everyday_todo_queryManyRow() []TODO {
	s := "select Id, Content, Done, Priority from everyday_todo_tbl"
	rows, err := db.Query(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil
	}
	defer rows.Close()

	var everyday_todos []TODO
	for rows.Next() {
		var everyday_todo TODO
		err := rows.Scan(&everyday_todo.Id, &everyday_todo.Content, &everyday_todo.Done, &everyday_todo.Priority)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil
		}
		everyday_todos = append(everyday_todos, everyday_todo)
	}
	return todos
}

func everyday_todo_deleteAll() {
	s := "delete from everyday_todo_tbl"
	_, err := db.Exec(s)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
