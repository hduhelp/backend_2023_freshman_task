package models

import (
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	Role     bool
	Email    string
	gorm.Model
}

type Todo struct {
	Username string
	ItemName string
	Detail   string
	EndTime  int64
	Finished bool
	gorm.Model
}

type StatusCode struct {
	Code        uint
	Description string
}

var DatabaseError = StatusCode{1, "DatabaseError"}
var UserExistsError = StatusCode{2, "UserExistsError"}
var UserNotExistsError = StatusCode{3, "UserNotExistsError"}
var ChangePWDSuccess = StatusCode{0, "ChangePWDSuccess"}
var ChangePWDFailed = StatusCode{4, "ChangePWDFailed"}
var AddTodoSuccess = StatusCode{0, "AddTodoSuccess"}
var ListTodoSuccess = StatusCode{0, "ListTodoSuccess"}
var ListTodoFailed = StatusCode{5, "ListTodoFailed"}
var DelTodoSuccess = StatusCode{0, "DelTodoSuccess"}
var LoginSuccess = StatusCode{0, "LoginSuccess"}
var LoginFailed = StatusCode{6, "LoginFailed"}
var RegisterSuccess = StatusCode{0, "RegisterSuccess"}
var ItemExistsError = StatusCode{7, "ItemExistsError"}
var ItemNotExistsError = StatusCode{8, "ItemNotExistsError"}
var ChangeTodoSuccess = StatusCode{0, "ChangeTodoSuccess"}
var ChangeTodoFailed = StatusCode{9, "ChangeTodoFailed"}
var EmailFormatError = StatusCode{10, "EmailFormatError"}
var SetEmailSuccess = StatusCode{0, "SetEmailSuccess"}
var SetEmailFailed = StatusCode{11, "SetEmailFailed"}
var GetUserSuccess = StatusCode{0, "GetUserSuccess"}
