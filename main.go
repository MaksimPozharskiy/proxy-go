package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(handleRequest),
	}

	serverMetrics := http.Server{
		Addr: ":8090",
		Handler: promhttp.Handler(),
	}

	go serverMetrics.ListenAndServe()

	log.Println("Starting proxy server on :8080 port")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}

func handleRequest(writer http.ResponseWriter, req *http.Request) {
	startReq := time.Now()

	targetUrl := os.Getenv("PROXY_TARGET_URL") + fmt.Sprintf("%s", req.URL)
	if targetUrl == "" {
		log.Fatal("PROXY_TARGET_URL in env not setted")
	}

	fmt.Println(req.Method, req.URL)
	
	proxyReq, err := http.NewRequest(req.Method, targetUrl, req.Body)
	if err != nil {
		http.Error(writer, "Error creaing proxy request", http.StatusInternalServerError)
		return
	}

	for name, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	customTransport := http.DefaultTransport
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Error sending proxy request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			writer.Header().Add(name, value)
		}
	}

	writer.WriteHeader(resp.StatusCode)

	io.Copy(writer, resp.Body)

	finishReq := time.Since(startReq)
	fmt.Printf("Request performs  %s\n", finishReq)
	requestsTotal.Inc()
}