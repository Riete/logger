package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

type Color string

const (
	// ANSI color, https://gist.github.com/iamnewton/8754917#file-bash-colors-md
	start  Color = "\033["
	end    Color = "\033[0m"
	Black  Color = "30m"
	Red    Color = "31m"
	Green  Color = "32m"
	Yellow Color = "33m"
	Blue   Color = "34m"
	Purple Color = "35m"
	Cyan   Color = "36m"
	White  Color = "37m"
	Gray   Color = "90m"
)

const (
	defaultCallerKey  = "source"
	defaultCallerSkip = 3
)

// DefaultColors overwrite or add additional level colors
var DefaultColors = map[slog.Level]Color{
	slog.LevelDebug: Gray,
	slog.LevelWarn:  Yellow,
	slog.LevelError: Red,
}

type colorWrapper struct {
	w     io.Writer
	level slog.Level
	lf    string
}

func (c *colorWrapper) Write(p []byte) (int, error) {
	color := DefaultColors[c.level]
	if color != "" {
		p = p[0 : len(p)-len(c.lf)]
		p = append([]byte(start+color), append(p, []byte(string(end)+c.lf)...)...)
	}
	return c.w.Write(p)
}

type Option func(*Logger)

func WithJSONFormat() Option {
	return func(l *Logger) {
		l.json = true
	}
}

func WithColor() Option {
	return func(l *Logger) {
		l.color = true
	}
}

func WithLogLevel(level slog.Level) Option {
	return func(l *Logger) {
		l.SetLevel(level)
	}
}

func WithMultiWriter(w io.Writer, others ...io.Writer) Option {
	return func(l *Logger) {
		l.w = io.MultiWriter(append(others, l.w, w)...)
		for _, other := range append(others, w) {
			if closer, ok := other.(io.Closer); ok {
				l.closers = append(l.closers, closer)
			}
		}
	}
}

func WithAttrs(attrs ...slog.Attr) Option {
	return func(l *Logger) {
		l.SetAttrs(attrs...)
	}
}

func WithCaller(key string, skip int) Option {
	return func(l *Logger) {
		l.caller = &caller{enable: true, key: key, skip: skip}
	}
}

func WithDefaultCaller() Option {
	return func(l *Logger) {
		l.caller = &caller{enable: true, key: defaultCallerKey, skip: defaultCallerSkip}
	}
}

type caller struct {
	enable bool
	key    string
	skip   int
}

func (c *caller) caller() string {
	_, file, line, _ := runtime.Caller(c.skip)
	return fmt.Sprintf("%s:%d", file, line)
}

type Logger struct {
	json    bool
	color   bool
	logger  *slog.Logger
	level   *slog.LevelVar
	closers []io.Closer
	w       io.Writer
	mu      sync.Mutex
	caller  *caller
}

func (l *Logger) SetLevel(level slog.Level) {
	l.level.Set(level)
}

func (l *Logger) SetAttrs(attrs ...slog.Attr) {
	for _, attr := range attrs {
		l.logger = l.logger.With(attr)
	}
}

func (l *Logger) Log(level slog.Level, msg string, args ...any) {
	if l.color {
		if wrapper, ok := l.w.(*colorWrapper); ok {
			l.mu.Lock()
			defer l.mu.Unlock()
			wrapper.level = level
		}
	}
	if l.caller.enable {
		args = append(args, l.caller.key, l.caller.caller())
	}
	l.logger.Log(context.Background(), level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.Log(slog.LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(msg string, args ...any) {
	l.Log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.Log(slog.LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.Log(slog.LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(msg string, args ...any) {
	l.Log(slog.LevelError, msg, args...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.Log(slog.LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(slog.LevelError, msg, args...)
	l.Close()
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Log(slog.LevelError, fmt.Sprintf(format, v...))
	l.Close()
	os.Exit(1)
}

func (l *Logger) Print(v ...any) {
	if len(v) == 0 {
		l.Log(slog.LevelInfo, "")
	}
	msg, ok := v[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", v[0])
	}
	l.Log(slog.LevelInfo, msg, v[1:]...)
}

func (l *Logger) Println(v ...any) {
	if len(v) == 0 {
		l.Log(slog.LevelInfo, "")
	}
	msg, ok := v[0].(string)
	if !ok {
		msg = fmt.Sprintf("%v", v[0])
	}
	l.Log(slog.LevelInfo, msg, v[1:]...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.Log(slog.LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Write(p []byte) (int, error) {
	return l.w.Write(p)
}

func (l *Logger) Close() error {
	for _, c := range l.closers {
		_ = c.Close()
	}
	return nil
}

func New(w io.Writer, options ...Option) *Logger {
	l := &Logger{level: new(slog.LevelVar), w: w}
	if closer, ok := w.(io.Closer); ok {
		l.closers = append(l.closers, closer)
	}
	for _, option := range options {
		option(l)
	}
	if l.color {
		if runtime.GOOS == "windows" {
			l.w = &colorWrapper{w: l.w, lf: "\r\n"}
		} else {
			l.w = &colorWrapper{w: l.w, lf: "\n"}
		}
	}
	var handler slog.Handler
	if l.json {
		handler = slog.NewJSONHandler(l.w, &slog.HandlerOptions{Level: l.level})
	} else {
		handler = slog.NewTextHandler(l.w, &slog.HandlerOptions{Level: l.level})
	}
	l.logger = slog.New(handler)
	return l
}
