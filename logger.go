package logger

import (
	"context"
	"io"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
	level  *slog.LevelVar
	w      io.Writer
}

func (l *Logger) SetLevel(level slog.Level) {
	l.level.Set(level)
}

func (l *Logger) SetAttrs(attrs ...any) {
	l.logger = l.logger.With(attrs...)
}

func (l *Logger) SetGroup(group string) {
	l.logger = l.logger.WithGroup(group)
}

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	l.logger.Log(context.Background(), level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, args...)
}

func (l *Logger) Write(p []byte) (int, error) {
	return l.w.Write(p)
}

// NewJSONLogger json log format, support write to multi writer
func NewJSONLogger(w io.Writer, others ...io.Writer) *Logger {
	if len(others) > 0 {
		w = io.MultiWriter(append(others, w)...)
	}
	level := new(slog.LevelVar)
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})),
		level:  level,
		w:      w,
	}
}

// NewTextLogger text log format, support write to multi writer
func NewTextLogger(w io.Writer, others ...io.Writer) *Logger {
	if len(others) > 0 {
		w = io.MultiWriter(append(others, w)...)
	}
	level := new(slog.LevelVar)
	return &Logger{
		logger: slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})),
		level:  level,
		w:      w,
	}
}
