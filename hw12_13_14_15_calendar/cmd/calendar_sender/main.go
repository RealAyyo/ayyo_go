package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/transport/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	eventQueue = "calendar_event"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("load config error: " + err.Error())
	}
	logg := logger.New(conf.Logger.Level)

	broker, err := rabbitmq.New(conf.RabbitMQ.User, conf.RabbitMQ.Password, conf.RabbitMQ.Host, conf.RabbitMQ.Port)
	if err != nil {
		logg.Error("failed to init broker: " + err.Error())
		os.Exit(1)
	}

	defer func(broker *rabbitmq.RabbitMQ) {
		err := broker.Close()
		if err != nil {
			logg.Error("failed to close broker: " + err.Error())
		}
	}(broker)

	var wg sync.WaitGroup
	wg.Add(1)
	logg.Info("Sender started")

	err = broker.QueueDeclare()
	if err != nil {
		logg.Error("error of queue declare: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Start consuming1")
	err = broker.Consume(eventQueue, handler)
	if err != nil {
		logg.Error("failed to consume topic: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("Start consuming")

	wg.Wait()
}

func handler(message amqp.Delivery) {
	log.Printf("Notify! %s", message.Body)
}
