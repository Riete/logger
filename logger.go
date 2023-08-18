package logger

import (
	"log/slog"

	"github.com/riete/logger/log"
	"github.com/riete/logger/writer"
)

type Logger interface {
	SetLevel(level slog.Level)
	Debug(msg string, v ...any)
	Info(msg string, v ...any)
	Warn(msg string, v ...any)
	Error(msg string, v ...any)
	Printf(format string, v ...any)
	Write(p []byte) (n int, err error)
}

func NewFileWriter(filename string, opts ...writer.FileWriterOption) writer.Writer {
	return writer.NewFileWriter(filename, opts...)
}

func NewLogger(level slog.Level, w writer.Writer) Logger {
	return log.NewLogger(level, w)
}

func NewJSONLogger(level slog.Level, w writer.Writer, attrs ...slog.Attr) Logger {
	return log.NewSLogger(level, true, w, attrs...)
}

func NewTextLogger(level slog.Level, w writer.Writer, attrs ...slog.Attr) Logger {
	return log.NewSLogger(level, false, w, attrs...)
}
