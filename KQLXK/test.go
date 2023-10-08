package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type logok struct {
	Username string
	Logok    string
}

func main() {
	db, _ := getDBConnection("logok1")
	db.Create(&logok{})
}
