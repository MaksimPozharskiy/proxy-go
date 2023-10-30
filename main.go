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
	signal.Notify(exit, os.		Interrupt, syscall.SIGTERM)

	go serverMetrics.ListenAndServe()

	log.Printf("Starting proxy server on %s port\n", os.Getenv("PROXY_SERVER_PORT"))

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Error starting proxy server: ", err)
		}
	}()

	<-exit	
	log.Println("Graceful shutdowning server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown error: ", err)
	}

}
