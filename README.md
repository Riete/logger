# logger

```
package main

import (
    "log/slog"
    "os"
    
    "github.com/riete/logger"
)

func main() {
    // create a file writer with DefaultFileRotator
    fw := logger.NewFileWriter("logs/test.log", logger.DefaultFileRotator)
    
    // create a network(udp or tcp) writer, 
    nw := logger.NewNetworkWriter("tcp", "127.0.0.1", 10*logger.SizeMiB)
    
    log := logger.New(
        fw,
        logger.WithJSONFormat(),                 // use json format
        logger.WithMultiWriter(os.Stdout, nw),   // log to multi target
        logger.WithColor(),                      // use ansi color
        logger.WithLogLevel(slog.LevelInfo),     // default log level, lower than it will be ignored
        logger.WithDefaultCaller(),              // source code position of the log statement with key "source" and skip 3 
        logger.WithCaller("source", 4),          // custom key and skip(more than 3 if wrapped) 
    )
    
    defer fw.Close() // flush buffered data 
    defer nw.Close() // close network connection
    or
    defer log.Close() // close all io.Closer
    
    log.Debug("debug log")
    log.Info("info log")
    log.Warn("warn log")
    log.Error("error log")
    log.Fatal("fatal log")
    
    log.Debugf("%s", "debug log")
    log.Infof("%s", "info log")
    log.Warnf("%s", "warn log")
    log.Errorf("%s", "error log")
    log.Fatalf("%s", "fatal log")

}

```