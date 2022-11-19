package main

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger   LoggerConf
	Queue    QueueConf
	Database DBConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level string
}

type QueueConf struct {
	Dsn  string
	Name string
}

type DBConf struct {
	Dsn string
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