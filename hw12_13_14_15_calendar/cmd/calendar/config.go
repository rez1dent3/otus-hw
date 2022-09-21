package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level string
}

func NewConfig(reader io.Reader) (*Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
