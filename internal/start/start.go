package start

import (
	"context"
	"log"
	"net/http"

	"github.com/lazy-electron-consulting/renogy-exporter/internal/config"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/metrics"
	"github.com/lazy-electron-consulting/renogy-exporter/internal/renogy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger = log.New(log.Writer(), "[start] ", log.Lmsgprefix|log.Flags())

// Run starts the server and runs until the given context is canceled
func Run(ctx context.Context, cfg *config.Config) error {
	logger.Println("starting")

	ctx, stop := context.WithCancel(ctx)
	defer stop()
	defer logger.Println("stopped")
	r, err := renogy.New(cfg.Modbus)
	if err != nil {
		return err
	}
	defer r.Close()
	err = metrics.Register(r, cfg.Gauges)
	if err != nil {
		return err
	}

	return runHttp(ctx, cfg.Address)
}

func runHttp(ctx context.Context, addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: addr}
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
