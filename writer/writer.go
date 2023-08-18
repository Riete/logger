package writer

import (
	"io"
)

type Writer interface {
	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	Rotate() bool
	Writer() io.Writer
}
