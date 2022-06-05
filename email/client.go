package email

import (
	"crypto/tls"
	"os"

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

	var password string
	if len(os.Getenv("YANDEX_PASSWORD")) == 0 {
		password = "cnxfedsnjsdrhdwl"
	} else {
		password = os.Getenv("YANDEX_PASSWORD")
	}

	d := gomail.NewDialer("smtp.yandex.com", 465, "sociumsocialmedia@yandex.com", password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)

	return err
}
