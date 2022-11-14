package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/exchanges/rabbitmq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yaml", "Path to configuration file")
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

	repo, err := storage.NewSchedulerStorage(ctx, config.Storage.Driver, config.Database.Dsn)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	defer func() {
		if err := repo.Close(); err != nil {
			logg.Error(err.Error())
		}
	}()

	amqp := rabbitmq.New(config.Queue.Dsn)
	sch := scheduler.New(repo, amqp, config.Queue.Name, logg, config.Scheduler.Notify, config.Scheduler.Cleaner)

	defer func() {
		err := amqp.Close()
		if err != nil {
			logg.Error(err.Error())
		}
	}()

	go func() {
		logg.Info("scheduler is running...")
		if err := sch.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
