package start

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/metrics"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/renogy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger = log.New(log.Writer(), "[start] ", log.Lmsgprefix|log.Flags())

// Config carries all the inputs needed to run the exporter
type Config struct {
	Port int
	Path string
}

// Run starts the server and runs until the given context is canceled
func Run(ctx context.Context, cfg *Config) error {
	logger.Println("starting")

	ctx, stop := context.WithCancel(ctx)
	defer stop()
	defer logger.Println("stopped")
	r, err := renogy.New(cfg.Path)
	if err != nil {
		return err
	}
	defer r.Close()
	err = metrics.Register(r)
	if err != nil {
		return err
	}

	return runHttp(ctx, cfg)
}

func runHttp(ctx context.Context, cfg *Config) error {
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	defer logger.Println("http server stopped")
	logger.Printf("http server started on %v\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Printf("error serving %v\n", err)
		return err
	}
	return nil
}
