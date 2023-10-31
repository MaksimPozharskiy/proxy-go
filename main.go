package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaksimPozharskiy/proxy-go/metrics"
	"github.com/MaksimPozharskiy/proxy-go/proxy"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go metrics.RunMetricsServer(ctx)
	go proxy.RunProxyServer(ctx)

	log.Printf("Starting proxy server on %s port\n", os.Getenv("PROXY_SERVER_PORT"))

	<-exit
	log.Println("Graceful shutdowning servers...")

	defer cancel()
}
