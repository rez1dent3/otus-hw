package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	logg.Debug("db driver: " + config.Storage.Driver)

	selStorage, err := storage.NewConnect(ctx, config.Storage.Driver, config.Database.Dsn)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	defer func() {
		if err := selStorage.Close(); err != nil {
			logg.Error(err.Error())
		}
	}()

	repo, ok := selStorage.(app.Storage)
	if !ok {
		logg.Error("an invalid object was received")
		return
	}

	calendar := app.New(logg, repo)

	serverGrpc := internalgrpc.NewServer(logg, calendar, config.Grpc.Host, config.Grpc.Port)
	serverHttp := internalhttp.NewServer(logg, calendar, config.Http.Host, config.Http.Port)

	logg.Info("calendar is running...")

	go func() {
		logg.Info("grpc is running...")
		if err := serverGrpc.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	go func() {
		logg.Info("http is running...")
		if err := serverHttp.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	<-ctx.Done()
}
