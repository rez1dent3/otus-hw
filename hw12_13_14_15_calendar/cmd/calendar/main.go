package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	config, err := NewConfig(file)
	if err != nil {
		log.Fatalln(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level, os.Stdout)
	var repo app.Storage
	if config.Storage.Driver == "postgres" {
		pgStorage := storage.NewPgStorage(config.Database.Dsn)
		err := pgStorage.Connect(ctx)
		if err != nil {
			logg.Error(err.Error())
			return
		}

		defer func(pgStorage *storage.PgStorage) {
			err := pgStorage.Close()
			if err != nil {
				logg.Error(err.Error())
			}
		}(pgStorage)

		repo = pgStorage
	} else {
		repo = storage.NewMemStorage()
	}

	logg.Debug("db driver: " + config.Storage.Driver)
	calendar := app.New(logg, repo)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

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
