package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger = log.New(log.Writer(), "[main] ", log.Lmsgprefix|log.Flags())

func main() {
	logger.Println("starting server")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT)
	defer stop()
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: ":8080"}
	defer srv.Close()
	go func() {
		defer logger.Println("server stopped")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("error serving %v\n", err)
			stop()
		}
	}()
	logger.Println("server started, Ctrl-C to exit")
	<-ctx.Done()
}
