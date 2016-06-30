package mailer

import (
	"../config"
	mailgun "github.com/mailgun/mailgun-go"
	"log"
)

var mailer mailgun.Mailgun

func Init() {
	mailer = mailgun.NewMailgun(config.Cfg.MailDomain, config.Cfg.MailPrivate, config.Cfg.MailPublic)
}

func sendMsg(from, subject, text, textHtml, to string) {
	// https://documentation.mailgun.com/api-sending.html#examples
	m := mailgun.NewMessage(from, subject, text, to)
	m.SetHtml(textHtml)

	if response, id, err := mailer.Send(m); err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent", id, response)
	}
}

func EmailConfirmation(email string) {

}

func RegistrationComplete(email string) {

}
