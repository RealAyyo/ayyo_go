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
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	StorageTypeSql    = "SQL"
	StorageTypeMemory = "MEMORY"
)

func main() {
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

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

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, conf.HTTP)

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
