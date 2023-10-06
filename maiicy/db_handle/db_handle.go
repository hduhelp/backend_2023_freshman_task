package db_handle

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"login-system/utils"
	"os"
	"time"
)

var db *sql.DB // 全局数据库对象

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

func ConnectDatabase(dbPath string) error {
	var isNotExist bool
	_, err := os.Stat(dbPath)
	isNotExist = os.IsNotExist(err)
	if isNotExist {
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}
		err2 := file.Close()
		if err2 != nil {
			return err2
		}
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if isNotExist {
		createTodoSQL := `CREATE TABLE todos (
							todo_id INTEGER PRIMARY KEY AUTOINCREMENT,
							user_id INTEGER NOT NULL,
							title VARCHAR(255) NOT NULL,
							completed BOOLEAN NOT NULL DEFAULT 0,
							created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
							due_date DATE);`

		createUserSQL := `CREATE TABLE IF NOT EXISTS users (
							id INTEGER PRIMARY KEY AUTOINCREMENT,
							username TEXT NOT NULL,
							password TEXT NOT NULL,
							info TEXT);`

		createJwtBlackListSQL := `CREATE TABLE jwt_blacklist (
								id INTEGER PRIMARY KEY AUTOINCREMENT,
								token TEXT NOT NULL,
								expiry TIMESTAMP NOT NULL);`

		_, err := db.Exec(createTodoSQL)
		if err != nil {
			return err
		}

		_, err = db.Exec(createUserSQL)
		if err != nil {
			return err
		}

		_, err = db.Exec(createJwtBlackListSQL)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertTodo(userID int, title string, dueDate string) error {
	// dueDate 以文本格式（YYYY-MM-DD）表示
	sqlStatement := `
		INSERT INTO todos (user_id, title, completed, due_date)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(sqlStatement, userID, title, false, dueDate)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTodo(todoID int) error {
	// 准备SQL语句，根据Todo项的ID删除记录
	sqlStatement := `
		DELETE FROM todos
		WHERE todo_id = ?
	`
	// 执行SQL语句
	_, err := db.Exec(sqlStatement, todoID)
	if err != nil {
		return err
	}
	return nil
}

func FindTodosByUserID(userID int) ([]Todo, error) {
	rows, err := db.Query("SELECT todo_id, user_id, title, completed, due_date, created_at FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed, &todo.DueDate, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func FindTodoByID(todoID int) (Todo, error) {
	var todo Todo

	query := "SELECT todo_id, user_id, title, completed, due_date, created_at FROM todos WHERE todo_id = ?"

	err := db.QueryRow(query, todoID).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed, &todo.DueDate, &todo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 如果没有找到匹配的Todo项，可以返回自定义错误或nil
			return todo, fmt.Errorf("todo项不存在")
		}
		// 处理其他查询错误
		return todo, err
	}

	// 返回找到的Todo项
	return todo, nil
}

func FindTodosBeforeTime(beforeTime time.Time, userID int) ([]Todo, error) {
	// 查询数据库中在指定时间之前、属于指定用户的Todo项
	rows, err := db.Query("SELECT todo_id, user_id, title, completed, due_date, created_at FROM todos WHERE due_date < ? AND user_id = ?", beforeTime, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// 处理关闭行时的错误
		}
	}(rows)

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed, &todo.DueDate, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func UpdateTodo(todoID int, newTitle string, completed bool, dueDate string) error {
	sqlStatement := `
		UPDATE todos
		SET title = ?, completed = ?, due_date = ?
		WHERE todo_id = ?
	`
	_, err := db.Exec(sqlStatement, newTitle, completed, dueDate, todoID)
	if err != nil {
		return err
	}
	return nil
}

func InsertJWTIntoBlacklist(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SecretKey), nil
	})

	if err != nil {
		return err
	}

	// 检查JWT是否已过期
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := time.Unix(int64(claims["exp"].(float64)), 0)

		sqlStatement := `
			INSERT INTO jwt_blacklist (token, expiry)
			VALUES (?, ?)
		`
		// 执行SQL语句
		_, err := db.Exec(sqlStatement, tokenString, expiry)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("JWT令牌无效或已过期")
}

func DeleteExpiredTokens() error {
	currentTime := time.Now()

	sqlStatement := `
		DELETE FROM jwt_blacklist
		WHERE expiry < ?
	`
	_, err := db.Exec(sqlStatement, currentTime)
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
	// 查询JWT令牌是否在黑名单表中
	rows, err := db.Query("SELECT COUNT(*) FROM jwt_blacklist WHERE token = ?", tokenString)
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var count int

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func InsertUser(username, password, info string) error {
	insertUserSQL := `INSERT INTO users (username, password, info) VALUES (?, ?, ?)`
	_, err := db.Exec(insertUserSQL, username, password, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(userID int) (User, error) {
	var user User
	query := "SELECT id, username, password, info FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.Info)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 如果没有找到匹配的用户，可以处理相关逻辑
			return User{}, fmt.Errorf("用户不存在")
		}
		return User{}, err
	}
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	var result User
	query := "SELECT id, username, password, info FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&result.ID, &result.Username, &result.Password, &result.Info)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf("用户不存在")
		}
		return User{}, err
	}
	return result, nil
}
