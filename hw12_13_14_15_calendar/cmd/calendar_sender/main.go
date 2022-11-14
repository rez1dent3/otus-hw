package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/consumers/sender"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/exchanges/rabbitmq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.yaml", "Path to configuration file")
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
	amqp := rabbitmq.New(config.Queue.Dsn)
	snd := sender.New(logg, amqp, config.Queue.Name)

	defer func() {
		err := amqp.Close()
		if err != nil {
			logg.Error(err.Error())
		}
	}()

	go func() {
		logg.Info("sender is running...")
		if err := snd.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
