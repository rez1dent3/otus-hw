package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
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
		log.Println(err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	config, err := NewConfig(file)
	if err != nil {
		log.Println(err)
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

	serverGrpc := internalgrpc.NewServer(logg, calendar, config.Grpc.Host, config.Grpc.Port)
	serverHttp := internalhttp.NewServer(logg, calendar, config.Http.Host, config.Http.Port)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverGrpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		err := serverGrpc.Start(ctx)
		if err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	if err := serverHttp.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
