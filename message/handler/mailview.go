package handler

import (
	"encoding/json"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"

	"github.com/nsqio/go-nsq"
)

type EmailService struct {
	ToEmail string
	Subject string
	Message string
}

type Test struct {
	Message string
}

type Consumer struct {
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (consumer *Consumer) SendEmailWithText(message *nsq.Message) error {
	var emailService EmailService
	if err := json.Unmarshal(message.Body, &emailService); err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("MAIL_FROM_ADDRESS"))
	mailer.SetHeader("To", emailService.ToEmail)
	mailer.SetHeader("Subject", emailService.Subject)
	mailer.SetBody("text/html", emailService.Message)

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

func (consumer *Consumer) Test(message *nsq.Message) error {
	var test Test
	if err := json.Unmarshal(message.Body, &test); err != nil {
		return err
	}

	log.Println("NSQ received message: ", test.Message)
	return nil
}
