package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/start"
)

var logger = log.New(log.Writer(), "[main] ", log.Lmsgprefix|log.Flags())

var (
	port = flag.Int("port", 8080, "what port to run the http server on")
	path = flag.String("path", "/dev/ttyUSB0", "where to read modbus data")
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT)
	defer stop()
	err := start.Run(ctx, &start.Config{
		Port: *port,
		Path: *path,
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatalf("exiting with errors %v\n", err)
	}
}
