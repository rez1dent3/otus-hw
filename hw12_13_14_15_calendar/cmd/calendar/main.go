package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
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
		return
	}

	defer func() {
		_ = file.Close()
	}()

	config, err := NewConfig(file)
	if err != nil {
		return
	}

	logg := logger.New(config.Logger.Level, os.Stdout)
	storage := memorystorage.New()
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
