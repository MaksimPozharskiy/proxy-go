package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaksimPozharskiy/proxy-go/metrics"
	"github.com/MaksimPozharskiy/proxy-go/proxy"
)

func main() {
	server := proxy.CreateProxyServer()

	serverMetrics := metrics.CreateMetricsServer()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := serverMetrics.ListenAndServe(); err != nil {
			log.Fatal("Error starting metrics proxy server: ", err)
		}
	}()

	log.Printf("Starting proxy server on %s port\n", os.Getenv("PROXY_SERVER_PORT"))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Error starting proxy server: ", err)
		}
	}()

	<-exit
	log.Println("Graceful shutdowning servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown server error: ", err)
	}

	if err := serverMetrics.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown metrics server error: ", err)
	}

}
