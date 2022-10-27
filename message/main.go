package main

import (
	"github.com/arvians-id/apriori/mail-service/handler"
	"github.com/arvians-id/apriori/mail-service/messaging"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	channelMail = "mail_channel"
	topicMail   = "mail_topic"

	channelStorage = "storage_channel"
	topicStorage   = "storage_topic"
)

func main() {
	handlerMessage := handler.NewConsumer()
	mailConsumer := messaging.ConsumerConfig{
		Topic:         topicMail,
		Channel:       channelMail,
		LookupAddress: "nsqlookupd:4161",
		MaxAttempts:   10,
		MaxInFlight:   100,
		Handler:       handlerMessage.SendEmailWithText,
	}

	mail := messaging.NewConsumer(mailConsumer)
	mail.Run()

	storageConsumer := messaging.ConsumerConfig{
		Topic:         topicStorage,
		Channel:       channelStorage,
		LookupAddress: "nsqlookupd:4161",
		MaxAttempts:   10,
		MaxInFlight:   100,
		Handler:       handlerMessage.Test,
	}

	storage := messaging.NewConsumer(storageConsumer)
	storage.Run()

	// keep app alive until terminated
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("Application terminated")
	}
}
