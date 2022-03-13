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

func main() {
	log.Println("starting server")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT)
	defer stop()
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: ":8080"}
	defer srv.Close()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stop()
		}
	}()
	log.Println("server started, Ctrl-C to exit")
	<-ctx.Done()
}
