package log

import (
	"log"
	"log/slog"
	"sync"

	"github.com/riete/logger/writer"
)

type Logger struct {
	w      writer.Writer
	logger *log.Logger
	mutex  sync.Mutex
	level  slog.Level
}

func (l *Logger) enable(level slog.Level) bool {
	return level >= l.level
}

func (l *Logger) checkSize() {
	if l.w.Rotate() {
		l.logger.SetOutput(l.w.Writer())
	}
}

func (l *Logger) log(prefix, msg string, v ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if prefix != "" {
		v = append([]any{prefix, msg}, v...)
	}
	l.logger.Println(v...)
	l.checkSize()
}

func (l *Logger) SetLevel(level slog.Level) {
	l.level = level
}

func (l *Logger) Debug(msg string, v ...any) {
	if l.enable(slog.LevelDebug) {
		l.log(slog.LevelDebug.String(), msg, v...)
	}
}

func (l *Logger) Info(msg string, v ...any) {
	if l.enable(slog.LevelInfo) {
		l.log(slog.LevelInfo.String(), msg, v...)
	}
}

func (l *Logger) Warn(msg string, v ...any) {
	if l.enable(slog.LevelWarn) {
		l.log(slog.LevelWarn.String(), msg, v...)
	}
}

func (l *Logger) Error(msg string, v ...any) {
	if l.enable(slog.LevelError) {
		l.log(slog.LevelError.String(), msg, v...)
	}
}

func (l *Logger) Write(p []byte) (int, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if n, err := l.w.Write(p); err != nil {
		return n, err
	} else {
		l.checkSize()
		return n, nil
	}
}

func (l *Logger) Printf(format string, v ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf(format, v...)
	l.checkSize()
}

func NewLogger(level slog.Level, w writer.Writer) *Logger {
	return &Logger{w: w, level: level, logger: log.New(w.Writer(), "", log.LstdFlags)}
}
