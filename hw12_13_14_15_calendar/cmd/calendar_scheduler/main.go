package main

import (
	"context"
	"log"
	"os"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/notifier"
	memorystorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/transport/rabbitmq"
)

const (
	StorageTypeSql    = "SQL"
	StorageTypeMemory = "MEMORY"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("load config error: " + err.Error())
	}
	logg := logger.New(conf.Logger.Level)

	var storage app.StorageService
	ctx := context.Background()

	switch conf.Storage.Type {
	case StorageTypeMemory:
		storage, err = memorystorage.New()
	case StorageTypeSql:
		storage, err = sqlstorage.New(ctx, conf.DB)
	}

	if err != nil {
		logg.Error("failed to init storage: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}

	defer func() {
		if closer, ok := storage.(sqlstorage.Closer); ok {
			err := closer.Close(ctx)
			if err != nil {
				log.Printf("Error closing storage: %v", err)
			}
		}
	}()

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

	err = broker.QueueDeclare()
	if err != nil {
		logg.Error("error of queue declare: " + err.Error())
		os.Exit(1)
	}

	calendar := app.New(logg, storage)
	go calendar.StartCleanUp()

	notifierInstance, err := notifier.New(broker, calendar, logg, conf.Scheduler.Interval)
	if err != nil {
		logg.Error("failed to init notifier: " + err.Error())
		os.Exit(1)
	}

	logg.Info("Scheduler started")
	notifierInstance.Start()
}
