package main

import (
	"log"
	"os"
	"proxy-go/metrics"
	"proxy-go/proxy"
)

func main() {
	server := proxy.CreateProxyServer()

	serverMetrics := metrics.CreateMetricsServer()

	go serverMetrics.ListenAndServe()

	log.Printf("Starting proxy server on %s port\n", os.Getenv("PROXY_SERVER_PORT"))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
