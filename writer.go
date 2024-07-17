package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
)

type FileWriter struct {
	f       *os.File
	rotator Rotator
	path    string
	wSize   *atomic.Int64
}

func (f *FileWriter) open() {
	f.f, _ = os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func (f *FileWriter) loadSize() {
	fs, _ := f.f.Stat()
	f.wSize.Store(fs.Size())
}

func (f *FileWriter) Writer() io.Writer {
	return f
}

func (f *FileWriter) Write(p []byte) (int, error) {
	n, err := f.f.Write(p)
	f.wSize.Add(int64(n))
	if f.wSize.Load() >= f.rotator.MaxSize() {
		_ = f.f.Close()
		f.rotator.Rotate(f.path)
		f.open()
		f.wSize.Store(0)
	}
	return n, err
}

func NewFileWriter(path string, rotator Rotator) *FileWriter {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	fw := &FileWriter{path: path, rotator: rotator, wSize: new(atomic.Int64)}
	fw.open()
	fw.loadSize()
	return fw
}
