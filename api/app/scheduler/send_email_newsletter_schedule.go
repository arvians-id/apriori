package scheduler

import (
	"apriori/service"
	"log"
)

type SendEmailNewsletterSchedule struct {
}

func (scheduler *SendEmailNewsletterSchedule) Run() {
	serviceEmail := service.NewEmailService()
	message := "Hello, this is a test email cron from apriori"
	err := serviceEmail.SendEmailWithText("widdy@gmail.com", "Test Scheduler", &message)
	if err != nil {
		log.Fatalln(err)
	}
}
