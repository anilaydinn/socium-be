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
	d := gomail.NewDialer("smtp.gmail.com", 587, "sociumsocialmedia@gmail.com", "sociumtestpassword")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)

	return err
}
