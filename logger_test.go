package logger

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	info := NewLogger("info.log", WithBase("loogs"), WithMaxRotate(10), WithMaxSizePerFile(50*SizeMiB))
	for {
		info.Info("a", "b", "c")
		info.Error("aa", "bb", "cc")
		info.Log("aaa", "bbb", "ccc")
		time.Sleep(time.Millisecond)
	}
}
