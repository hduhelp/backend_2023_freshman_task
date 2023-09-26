package models

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(to string, content string) {
	em := email.NewEmail()
	em.From = fmt.Sprintf("todo <%s>", Env.GetString("email.sender"))
	fmt.Println(em.From)
	em.To = []string{to}
	em.Subject = "Todo"
	em.Text = []byte(content)
	server := fmt.Sprintf("%s:%s", Env.GetString("email.server"), Env.GetString("email.port"))
	err := em.Send(server, smtp.PlainAuth("", Env.GetString("email.sender"), Env.GetString("email.secret"), Env.GetString("email.server")))
	if err != nil {
		Logger.Error(err)
	} else {
		Logger.Info(fmt.Sprintf("Send to %s Successfully", to))
	}
}
