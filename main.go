package main

import (
	"log"
)

func main() {
	server := createProxyServer()

	serverMetrics := createMetricsServer()

	go serverMetrics.ListenAndServe()

	log.Println("Starting proxy server on :8080 port")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
