package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lazy-electron-consulting/renology-exporter/internal/start"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT)
	defer stop()
	start.Run(ctx, &start.Config{
		Port: 8080,
	})
}
