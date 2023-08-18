package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	SizeMiB               = 1 << 20
	defaultMaxSizePerFile = 20 * SizeMiB
	defaultMaxRotate      = 5
	defaultBase           = "logs"
)

type FileWriterOption func(*FileWriter)

func WithBase(base string) FileWriterOption {
	return func(w *FileWriter) {
		w.base = base
	}
}

func WithMaxSizePerFile(max int64) FileWriterOption {
	return func(w *FileWriter) {
		w.maxSizePerFile = max
	}
}

func WithMaxRotate(max int64) FileWriterOption {
	return func(w *FileWriter) {
		w.maxRotate = max
	}
}

type FileWriter struct {
	base           string
	filename       string
	f              *os.File
	maxSizePerFile int64
	maxRotate      int64
}

func (f *FileWriter) fullFilename() string {
	return filepath.Join(f.base, f.filename)
}

func (f *FileWriter) rotateFilename(i int64) string {
	return fmt.Sprintf("%s.%d", f.fullFilename(), i)
}

func (f *FileWriter) newWriter() {
	f.f, _ = os.OpenFile(f.fullFilename(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func (f *FileWriter) Rotate() bool {
	s, _ := f.f.Stat()
	rotateNeeded := s.Size() >= f.maxSizePerFile
	if rotateNeeded {
		_ = f.f.Close()
		for i := f.maxRotate; i > 0; i-- {
			src := f.rotateFilename(i - 1)
			dest := f.rotateFilename(i)
			if i == 1 {
				src = f.fullFilename()
			}
			_ = os.Rename(src, dest)
		}
		f.newWriter()
	}
	return rotateNeeded
}

func (f *FileWriter) Writer() io.Writer {
	return f.f
}

func (f *FileWriter) Write(b []byte) (int, error) {
	return f.f.Write(b)
}

func (f *FileWriter) WriteString(s string) (int, error) {
	return f.f.WriteString(s)
}

func NewFileWriter(filename string, opts ...FileWriterOption) *FileWriter {
	w := &FileWriter{
		base:           defaultBase,
		filename:       filename,
		maxSizePerFile: defaultMaxSizePerFile,
		maxRotate:      defaultMaxRotate,
	}
	for _, opt := range opts {
		opt(w)
	}
	_ = os.MkdirAll(w.base, 0755)
	w.newWriter()
	return w
}
