package main

import (
	"io"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConf
	Queue  QueueConf
}

type LoggerConf struct {
	Level string
}

type QueueConf struct {
	Dsn  string
	Name string
}

func NewConfig(reader io.Reader) (*Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
