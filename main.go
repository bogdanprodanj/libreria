package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/libreria/config"
	"github.com/libreria/server/http"
	"github.com/libreria/server/http/handlers"
	"github.com/libreria/service/book"
	"github.com/libreria/storage/postgres"
	log "github.com/sirupsen/logrus"
)

func main() { // nolint:funlen
	// read service cfg from os env
	cfg, err := config.New()
	if err != nil {
		log.WithError(err).Fatal("config init error")
	}
	// init logger
	initLogger(cfg.LogLevel)
	log.Info("service starting...")

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// create DB connection
	pg, err := postgres.New(ctx, wg, cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("postgres init error")
	}

	// create service
	bookSrv := book.New(pg)

	// initializing http server
	httpSrv := http.New(
		cfg.HTTPServer,
		handlers.New(bookSrv),
	)
	// run srv
	httpSrv.Run(ctx, wg)

	log.Info("app is running now")

	// wait while services work
	wg.Wait()
	log.Info("service stopped")
}

func initLogger(logLevel string) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
