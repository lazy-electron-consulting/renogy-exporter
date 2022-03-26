package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/config"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/start"
)

var logger = log.New(log.Writer(), "[main] ", log.Lmsgprefix|log.Flags())

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] config-file\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.ReadYaml(flag.Arg(0))
	if err != nil {
		logger.Fatalf("could not read config %s: %v", flag.Arg(0), err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT)
	defer stop()
	err = start.Run(ctx, cfg)
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatalf("exiting with errors %v\n", err)
	}
}
