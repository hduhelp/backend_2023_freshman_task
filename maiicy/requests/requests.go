package requests

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Info     string `json:"info"`
}

type AddTodoRequest struct {
	Title string `json:"title" binding:"required"`
	Date  string `json:"date" binding:"required"`
}

type UpdateTodoRequest struct {
	TodoID    int    `json:"todo_id"  binding:"required"`
	Title     string `json:"title"`
	Date      string `json:"date"`
	Completed bool   `json:"completed"`
}
