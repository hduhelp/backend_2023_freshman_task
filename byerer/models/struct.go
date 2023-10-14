package models

type User struct {
	UserID   int64  `gorm:"primaryKey;column:userID" json:"userID"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

type TODO struct {
	ID      string `gorm:"primaryKey;column:id" json:"id"` //update
	UserID  int64  `gorm:"column:userID" json:"userID"`
	Content string `gorm:"column:content" json:"content"`
	Done    bool   `gorm:"column:done" json:"done"`
}
