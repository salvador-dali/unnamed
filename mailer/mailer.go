package mailer

import (
	"../config"
	"fmt"
	mailgun "github.com/mailgun/mailgun-go"
	"log"
)

var mailer mailgun.Mailgun

const (
	emailFrom = "registration@unnamed.com"
)

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

// getEmail returns either an email address provided to it, or an PROJ_TEST_EMAIL if PROJ_IS_TEST
func getEmail(email string) string {
	if config.Cfg.IsTest {
		return config.Cfg.TestEmail
	}

	return email
}

func EmailConfirmation(email, code string) {
	email = getEmail(email)
	text := fmt.Sprintf("Your confirmation code is: %s", code)
	textHtml := fmt.Sprintf("Your confirmation code is: <b>%s</b>", code)
	sendMsg(emailFrom, "Please confirm your registration", text, textHtml, email)
}
