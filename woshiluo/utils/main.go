//
// main.go
// Copyright (C) 2023 Woshiluo Luo <woshiluo.luo@outlook.com>
//
// Distributed under terms of the GNU AGPLv3+ license.
//

package utils

import (
	"fmt"
	"time"

	gomail "gopkg.in/mail.v2"

	log "github.com/sirupsen/logrus"

	"todolist/models"
)

type Email struct {
	Address     string `toml:"address"`
	SmtpUser    string `toml:"smtp_user"`
	SmtpAddress string `toml:"smtp_address"`
	SmtpPort    uint32 `toml:"smtp_port"`
	Password    string `toml:"password"`
}

type Config struct {
	DatabaseFile string `toml:"database_file"`
	Email        Email  `toml:"email"`
}

var ConfigData Config

func update_token() {
	log.Trace("Try update token")
	var tokens []models.Token

	if err := models.Db.Find(&tokens).Error; err != nil {
		log.Error(err)
		return
	}

	for _, token := range tokens {
		if time.Now().Sub(token.UpdatedAt) >= time.Minute*24*60 {
			log.Info("Token ", token, " is out of date, deleted")
			models.Db.Delete(&token)
		} else {
			log.Trace("Token ", token, " is still in use")
		}
	}
}

func send_mail(todo models.Todo) {
	log.Trace("Try send mail")

	var user models.User
	if err := models.Db.First(&user, todo.UserID).Error; err != nil {
		log.Error(err)
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", ConfigData.Email.Address)
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Todo notify")
	mail.SetBody("text/plain", fmt.Sprintln("Your todo ", todo, " is still undone."))

	dail := gomail.NewDialer(ConfigData.Email.SmtpAddress, int(ConfigData.Email.SmtpPort), ConfigData.Email.SmtpUser, ConfigData.Email.Password)

	if err := dail.DialAndSend(mail); err != nil {
		log.Error("Failed to send mail: ", err)
	}

	log.Trace("Mail Sent Successfully")
}

func notify_todo() {
	log.Trace("Try notify todo")
	var todos []models.Todo

	if err := models.Db.Find(&todos).Error; err != nil {
		log.Error(err)
		return
	}

	current_time := time.Now()

	for _, todo := range todos {
		if todo.DueDate == 0 || todo.Done {
			continue
		}
		var diff = time.Unix(int64(todo.DueDate), 0).Sub(current_time)

		// Cause we only run this function one time in every 30 seconds.
		// So it should be only run one time.
		// TODO: Rewrite in a more elegant way
		if diff <= time.Hour*24 && diff > (time.Hour*24-time.Second*30) {
			log.Info("Token ", todo, " need to nitify")
			send_mail(todo)
		} else {
			log.Trace("Token ", todo, " dose not need to notify", diff)
		}
	}
}

func Monitor() {
	for true {
		update_token()
		notify_todo()
		time.Sleep(30 * time.Second)
	}
}
