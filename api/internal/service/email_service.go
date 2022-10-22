package service

import (
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type EmailServiceImpl struct {
}

func NewEmailService() EmailService {
	return &EmailServiceImpl{}
}

func (service *EmailServiceImpl) SendEmailWithText(toEmail string, subject string, message *string) error {
	go func() {
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", os.Getenv("MAIL_FROM_ADDRESS"))
		mailer.SetHeader("To", toEmail)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", *message)

		port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
		if err != nil {
			log.Fatal(err)
		}
		dialer := gomail.NewDialer(
			os.Getenv("MAIL_HOST"),
			port,
			os.Getenv("MAIL_USERNAME"),
			os.Getenv("MAIL_PASSWORD"),
		)

		err = dialer.DialAndSend(mailer)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
