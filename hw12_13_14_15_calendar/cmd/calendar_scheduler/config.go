package main

import (
	"io"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger    LoggerConf
	Database  DBConf
	Scheduler SchedulerConf
	Queue     QueueConf
	Storage   StorageConf
}

type QueueConf struct {
	Dsn  string
	Name string
}

type LoggerConf struct {
	Level string
}

type DBConf struct {
	Dsn string
}

type SchedulerConf struct {
	Notify  time.Duration
	Cleaner time.Duration
}

type StorageConf struct {
	Driver string
}

func NewConfig(reader io.Reader) (*Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
