package main

import (
	"log"
	"os"
)

func main() {
	server := createProxyServer()

	serverMetrics := createMetricsServer()

	go serverMetrics.ListenAndServe()

	log.Printf("Starting proxy server on %s port\n", os.Getenv("PROXY_SERVER_PORT"))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}
