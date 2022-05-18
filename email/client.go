package email

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

const mail = "sociumsocialmedia@gmail.com"

type Client interface {
	SendMail(registeredUser, message string) error
}

func SendMail(registeredUser, subject, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", mail)
	m.SetHeader("To", registeredUser)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", message)
	d := gomail.NewDialer("smtp-relay.sendinblue.com", 587, "anlaydinn22@gmail.com", "zU7Btjv9IYwM8dny")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)

	return err
}
