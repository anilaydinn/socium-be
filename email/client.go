package email

import (
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

const mail = "sociumsocialmedia@yandex.com"

type Client interface {
	SendMail(registeredUser, message string) error
}

func SendMail(registeredUser, subject, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", mail)
	m.SetHeader("To", registeredUser)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", message)
	d := gomail.NewDialer("smtp.yandex.com", 465, "sociumsocialmedia@yandex.com", os.GetEnv("YANDEX_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)

	return err
}
