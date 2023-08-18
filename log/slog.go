package log

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/riete/logger/writer"
)

type SLogger struct {
	json   bool
	attrs  []slog.Attr
	w      writer.Writer
	logger *slog.Logger
	level  *slog.LevelVar
	mutex  sync.Mutex
}

func (s *SLogger) checkSize() {
	if s.w.Rotate() {
		s.newLogger()
	}
}

func (s *SLogger) newLogger() {
	if s.json {
		s.logger = slog.New(slog.NewJSONHandler(s.w.Writer(), &slog.HandlerOptions{Level: s.level}).WithAttrs(s.attrs))
	} else {
		s.logger = slog.New(slog.NewTextHandler(s.w.Writer(), &slog.HandlerOptions{Level: s.level}).WithAttrs(s.attrs))
	}
}

func (s *SLogger) log(level slog.Level, msg string, v ...any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.logger.Log(context.Background(), level, msg, v...)
	s.checkSize()
}

func (s *SLogger) SetLevel(level slog.Level) {
	s.level.Set(level)
}

func (s *SLogger) Debug(msg string, v ...any) {
	s.log(slog.LevelDebug, msg, v...)
}

func (s *SLogger) Info(msg string, v ...any) {
	s.log(slog.LevelInfo, msg, v...)
}

func (s *SLogger) Warn(msg string, v ...any) {
	s.log(slog.LevelWarn, msg, v...)
}

func (s *SLogger) Error(msg string, v ...any) {
	s.log(slog.LevelError, msg, v...)
}

func (s *SLogger) Write(p []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if n, err := s.w.Write(p); err != nil {
		return n, err
	} else {
		s.checkSize()
		return n, nil
	}
}

func (s *SLogger) Printf(format string, v ...any) {
	s.mutex.Lock()
	defer s.mutex.Lock()
	if _, err := s.w.WriteString(fmt.Sprintf(format, v...)); err == nil {
		s.checkSize()
	}
}

func NewSLogger(level slog.Level, json bool, w writer.Writer, attrs ...slog.Attr) *SLogger {
	s := &SLogger{
		json:  json,
		attrs: attrs,
		w:     w,
		level: &slog.LevelVar{},
	}
	s.SetLevel(level)
	s.newLogger()
	return s
}
