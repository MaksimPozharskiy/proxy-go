package main

import (
	"log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounter(
			prometheus.CounterOpts{
					Name: "proxy_go_requests_total",
					Help: "Total number of requests",
			},
	)
)

func init() { 
	prometheus.MustRegister(requestsTotal)
}

func main() {
	server := startProxyServer()

	serverMetrics := startMetricsServer()

	go serverMetrics.ListenAndServe()

	log.Println("Starting proxy server on :8080 port")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
