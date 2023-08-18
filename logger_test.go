package logger

import (
	"log/slog"
	"testing"
)

var w = NewFileWriter("test.log")

func TestNewLogger(t *testing.T) {
	l := NewLogger(slog.LevelInfo, w)
	for {
		l.Debug("debug")
		l.Info("info")
		l.Warn("warn")
		l.Error("error")
		_, _ = l.Write([]byte("xxxxxxxxxxxxx\n"))
	}
}

func TestNewJSONLogger(t *testing.T) {
	l := NewJSONLogger(slog.LevelInfo, w)
	for {
		l.Debug("debug")
		l.Info("info")
		l.Warn("warn")
		l.Error("error")
		_, _ = l.Write([]byte("xxxxxxxxxxxxx\n"))
	}
}

func TestNewTextLogger(t *testing.T) {
	l := NewTextLogger(slog.LevelInfo, w)
	for {
		l.Debug("debug")
		l.Info("info")
		l.Warn("warn")
		l.Error("error")
		_, _ = l.Write([]byte("xxxxxxxxxxxxx\n"))
	}
}
