package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	config2 "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/server/http"
	storage2 "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config2.NewConfig()

	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	var err error
	ctx := context.Background()

	switch config.Storage.Type {
	case "MEMORY":
		storage, err = memorystorage.New()
	case "SQL":
		storage, err = sqlstorage.New(ctx, config.Db)
	}

	err = storage.AddEvent(ctx, &storage2.Event{
		ID:       1,
		UserId:   2,
		Title:    "hello",
		Date:     time.Now(),
		Duration: "5h10m0s",
	})

	if err != nil {
		panic(err)
	}

	storage.ListEvents(ctx, 2, time.Now(), time.Now().Add(time.Hour*10))

	defer func() {
		if closer, ok := storage.(sqlstorage.Closer); ok {
			err := closer.Close(ctx)
			if err != nil {
				log.Printf("Error closing storage: %v", err)
			}
		}
	}()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
