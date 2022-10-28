package handler

import (
	"encoding/json"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"

	"github.com/nsqio/go-nsq"
)

type MailValue struct {
	ToEmail string
	Subject string
	Message string
}

type MailService struct {
}

func NewMailService() *MailService {
	return &MailService{}
}

func (consumer *MailService) SendEmailWithText(message *nsq.Message) error {
	var mailValue MailValue
	if err := json.Unmarshal(message.Body, &mailValue); err != nil {
		log.Println("[EmailService][Unmarshal] unable to unmarshal data, err: ", err.Error())
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("MAIL_FROM_ADDRESS"))
	mailer.SetHeader("To", mailValue.ToEmail)
	mailer.SetHeader("Subject", mailValue.Subject)
	mailer.SetBody("text/html", mailValue.Message)

	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Println("[EmailService][SendEmailWithText] problem in conversion string to integer, err: ", err.Error())
		return err
	}
	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		port,
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Println("[EmailService][SendEmailWithText] problem in mail sender, err: ", err.Error())
		return err
	}

	return nil
}
