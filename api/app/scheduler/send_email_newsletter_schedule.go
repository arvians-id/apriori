package scheduler

import (
	"apriori/service"
	"log"
)

type SendEmailNewsletterSchedule struct {
}

func (scheduler *SendEmailNewsletterSchedule) Run() {
	serviceEmail := service.NewEmailService()
	err := serviceEmail.SendEmailWithText("widdy@gmail.com", "Hallo bang ini cron")
	if err != nil {
		log.Fatalln(err)
	}
}
