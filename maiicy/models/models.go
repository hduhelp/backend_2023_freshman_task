package models

type User struct {
	ID       int
	Username string
	Password string
	Info     string
}

type Todo struct {
	ID        int
	UserID    int
	Title     string
	Completed bool
	DueDate   string
	CreatedAt string
}
