package db_handle

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"login-system/models"
	"login-system/utils"
	"time"
)

var db *gorm.DB // 全局数据库对象

func ConnectDatabase(dbPath string) error {
	var err error

	// 连接到 SQLite 数据库
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	err = db.AutoMigrate(&models.Todo{}, &models.User{}, &models.JWTBlacklist{})
	if err != nil {
		return fmt.Errorf("数据库创建表失败: %w", err)
	}

	return nil
}

func InsertTodo(userID uint, title string, dueDate string) error {
	parsedDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return fmt.Errorf("日期解析失败: %w", err)
	}

	newTodo := models.Todo{
		UserID:    uint(userID),
		Title:     title,
		Completed: false,
		DueDate:   parsedDueDate,
	}

	result := db.Create(&newTodo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteTodo(todoID uint) error {

	var todo models.Todo
	result := db.First(&todo, todoID)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&todo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FindTodosByUserID(userID uint) ([]models.Todo, error) {
	var todos []models.Todo

	// 使用 GORM 查询符合条件的记录
	result := db.Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}

	return todos, nil
}

func FindTodoByID(todoID uint) (models.Todo, error) {
	var todo models.Todo

	// 使用 GORM 查询符合条件的记录
	result := db.First(&todo, todoID)
	if result.Error != nil {
		return todo, result.Error
	}

	return todo, nil
}

func FindTodosBeforeTime(beforeTime time.Time, userID uint) ([]models.Todo, error) {
	var todos []models.Todo

	result := db.Where("due_date < ? AND user_id = ?", beforeTime, userID).Find(&todos)
	if result.Error != nil {
		return nil, result.Error
	}

	return todos, nil
}

func UpdateTodo(todoID uint, newTitle string, completed bool, dueDate string) error {
	var todo models.Todo

	// 查找要更新的记录
	result := db.First(&todo, todoID)
	if result.Error != nil {
		return result.Error
	}

	// 更新记录
	todo.Title = newTitle
	todo.Completed = completed
	// 解析日期字符串
	parsedDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return fmt.Errorf("日期解析失败: %w", err)
	}
	todo.DueDate = parsedDueDate

	result = db.Save(&todo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func InsertJWTIntoBlacklist(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.SecretKey), nil
	})

	if err != nil {
		return fmt.Errorf("JWT 令牌解析失败: %w", err)
	}

	// 检查JWT是否已过期
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := time.Unix(int64(claims["exp"].(float64)), 0)

		jwtBlacklist := models.JWTBlacklist{
			Token:  tokenString,
			Expiry: expiry,
		}

		// 使用 GORM 创建 JWT 黑名单记录
		result := db.Create(&jwtBlacklist)
		if result.Error != nil {
			return result.Error
		}

		return nil
	}

	return fmt.Errorf("JWT令牌无效或已过期")
}

func DeleteExpiredTokens() error {
	currentTime := time.Now()

	// 使用 GORM 删除过期的令牌记录
	result := db.Where("expiry < ?", currentTime).Delete(&models.JWTBlacklist{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func IsTokenBlacklisted(tokenString string) (bool, error) {
	// 使用 GORM 查询JWT令牌是否在黑名单表中
	var count int64
	result := db.Model(&models.JWTBlacklist{}).Where("token = ?", tokenString).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func InsertUser(username, password, info string) error {
	// 创建 User 记录
	newUser := models.User{
		Username: username,
		Password: password,
		Info:     info,
	}

	// 使用 GORM 创建记录
	result := db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserByID(userID uint) (models.User, error) {
	var user models.User

	// 使用 GORM 查询符合条件的记录
	result := db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("用户不存在")
		}
		return models.User{}, result.Error
	}

	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var userResult models.User

	// 使用 GORM 查询符合条件的记录
	result := db.Where("username = ?", username).First(&userResult)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, fmt.Errorf("用户不存在")
		}
		return models.User{}, result.Error
	}

	return userResult, nil
}
