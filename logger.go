package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type Logger struct {
	attrs  []any
	logger *slog.Logger
	level  *slog.LevelVar
	w      io.Writer
}

func (l *Logger) SetLevel(level slog.Level) {
	l.level.Set(level)
}

func (l *Logger) ClearAttrs() {
	l.attrs = []any{}
}

func (l *Logger) SetAttrs(attrs ...slog.Attr) {
	for _, attr := range attrs {
		l.attrs = append(l.attrs, attr)
	}
}

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	l.logger.Log(context.Background(), level, msg, append(args, l.attrs...)...)
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

func (l *Logger) Print(v ...any) {
	if len(v) == 0 {
		l.Info("")
	}
	msg, ok := v[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", v[0])
	}
	l.Info(msg, v[1:]...)
}

func (l *Logger) Println(v ...any) {
	l.Print(v...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.Info(fmt.Sprintf(format, v...))
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
