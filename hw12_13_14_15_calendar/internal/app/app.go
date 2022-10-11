package app

import (
	"context"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
}

type Storage interface {
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
