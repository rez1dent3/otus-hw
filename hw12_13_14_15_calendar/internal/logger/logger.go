package logger

import (
	"io"
)

type levelType uint8

const (
	levelOff levelType = iota
	levelError
	levelWarning
	levelInfo
	levelDebug
)

type Logger struct {
	writer io.Writer
	level  levelType
}

func New(level string, writer io.Writer) *Logger {
	var uintLevel levelType
	switch level {
	case "error":
		uintLevel = levelError
	case "warning":
		uintLevel = levelWarning
	case "info":
		uintLevel = levelInfo
	case "debug":
		uintLevel = levelDebug
	default:
		uintLevel = levelOff
	}
	return &Logger{level: uintLevel, writer: writer}
}

func (l *Logger) Error(msg string) {
	l.log(levelError, msg)
}

func (l *Logger) Warning(msg string) {
	l.log(levelWarning, msg)
}

func (l *Logger) Info(msg string) {
	l.log(levelInfo, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(levelDebug, msg)
}

func (l *Logger) log(level levelType, msg string) {
	if l.level >= level {
		_, _ = l.writer.Write([]byte(msg + "\n"))
	}
}
