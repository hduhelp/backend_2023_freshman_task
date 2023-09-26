package models

import (
	"encoding/base64"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func InitCron() *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	if _, err := c.AddFunc("* * * * * *", do); err != nil {
		Logger.Fatalln(err)
	}
	return c
}

func getUndo() (StatusCode, []Todo) {
	var todos []Todo
	res := db.Model(&Todo{}).Not(Todo{Finished: true}).Find(&todos)
	if res.Error != nil {
		Logger.Warning(res.Error)
		return ListTodoFailed, []Todo{}
	}
	for k, todo := range todos {
		decoded, _ := base64.StdEncoding.DecodeString(todo.Detail)
		todos[k].Detail = string(decoded)
	}
	return ListTodoSuccess, todos
}

func do() {
	status, todos := getUndo()
	if status != ListTodoSuccess {
		Logger.Warning("cron:", status.Description)
	} else {
		for _, todo := range todos {
			if time.Now().Unix() >= todo.EndTime-1800 {
				status, user := GetUserByName(todo.Username)
				if status == GetUserSuccess {
					if user.Email == "" {
						Logger.Warning("%s unset email address")
						continue
					}
					content := fmt.Sprintf("Todo %s 马上就要到期了：\n%s", todo.ItemName, todo.Detail)
					content = fmt.Sprintf("%s\n\n截至时间：%s", content, time.Unix(todo.EndTime, 0).Format("2006-01-02 03:04:05 PM"))
					SendEmail(user.Email, content)
					if status, _ := Done(todo.ItemName, user.Username); status != ChangeTodoSuccess {
						Logger.Warning("Do cron %s from %s Error", todo.ItemName, user.Username)
						continue
					} else {
						Logger.Info("Do cron %s from %s Success", todo.ItemName, user.Username)
						continue
					}
				}
				Logger.Warning(status.Description)
			}
		}
	}
}
