package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	SizeMiB               = 1 << 20
	defaultMaxSizePerFile = 20 * SizeMiB
	defaultMaxRotate      = 5
	defaultBase           = "logs"
)

type Logger struct {
	base           string
	filename       string
	f              *os.File
	logger         *log.Logger
	maxSizePerFile int64
	maxRotate      int64
	mutex          sync.Mutex
}

type Option func(*Logger)

func WithBase(base string) Option {
	return func(l *Logger) {
		l.base = base
	}
}

func WithMaxSizePerFile(max int64) Option {
	return func(l *Logger) {
		l.maxSizePerFile = max
	}
}

func WithMaxRotate(max int64) Option {
	return func(l *Logger) {
		l.maxRotate = max
	}
}

func (l *Logger) fullFilename() string {
	return filepath.Join(l.base, l.filename)
}

func (l *Logger) rotateFilename(i int64) string {
	return fmt.Sprintf("%s.%d", l.fullFilename(), i)
}

func (l *Logger) openFile() {
	l.f, _ = os.OpenFile(l.fullFilename(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func (l *Logger) rotate() {
	for i := l.maxRotate; i > 0; i-- {
		src := l.rotateFilename(i - 1)
		dest := l.rotateFilename(i)
		if i == 1 {
			src = l.fullFilename()
		}
		_ = os.Rename(src, dest)
	}
	l.openFile()
	l.logger.SetOutput(l.f)
}

func (l *Logger) write(prefix string, v ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if prefix != "" {
		v = append([]any{prefix}, v...)
	}
	l.logger.Println(v...)
	s, _ := l.f.Stat()
	if s.Size() >= l.maxSizePerFile {
		_ = l.f.Close()
		l.rotate()
	}
}

func (l *Logger) Printf(format string, v ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logger.Printf(format, v...)
	s, _ := l.f.Stat()
	if s.Size() >= l.maxSizePerFile {
		_ = l.f.Close()
		l.rotate()
	}
}

func (l *Logger) Info(v ...any) {
	l.write("[INFO]", v...)
}

func (l *Logger) Error(v ...any) {
	l.write("[ERROR]", v...)
}

func (l *Logger) Log(v ...any) {
	l.write("", v...)
}

func (l *Logger) Writer() io.Writer {
	return l.logger.Writer()
}

func NewLogger(filename string, options ...Option) *Logger {
	l := &Logger{
		base:           defaultBase,
		filename:       filename,
		maxSizePerFile: defaultMaxSizePerFile,
		maxRotate:      defaultMaxRotate,
	}
	for _, option := range options {
		option(l)
	}
	_ = os.MkdirAll(l.base, 0755)
	l.openFile()
	l.logger = log.New(l.f, "", log.LstdFlags)
	return l
}
