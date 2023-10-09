package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

var proxyTargetUrl = os.Getenv("PROXY_TARGET_URL")

func createProxyServer() *http.Server {
	http.HandleFunc("/", handleRequest)

	server := &http.Server{
		Addr: ":" + os.Getenv("PROXY_SERVER_PORT"),
	}

	return server
}

func handleRequest(writer http.ResponseWriter, req *http.Request) {
	startReq := time.Now()

	targetUrl := proxyTargetUrl + fmt.Sprintf("%s", req.URL)
	if targetUrl == "" {
		log.Fatal("PROXY_TARGET_URL in env not setted")
	}

	fmt.Println(req.Method, req.URL)

	proxyReq, err := http.NewRequest(req.Method, targetUrl, req.Body)
	if err != nil {
		http.Error(writer, "Error creaing proxy request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = req.Header.Clone()

	resp, err := getRequestWithRetry(proxyReq)

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

	finishReq := time.Since(startReq).Seconds()
	throughput := float64(resp.ContentLength) / finishReq

	fmt.Printf("Request performs  %vs.\n", finishReq)
	responseTimeHistogram.Observe(finishReq)
	requestsTotal.Inc()
	httpStatusCount.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
	throughputHistogram.Observe(throughput)
}

func getRequest(req *http.Request) (*http.Response, error) {
	customTransport := http.DefaultTransport
	resp, err := customTransport.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getRequestWithRetry(req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response

	for _, backoff := range backoffSchedule {
		resp, err = getRequest(req)

		if err == nil {
			return resp, nil
		}

		fmt.Printf("Request error: %v\n", err)
		fmt.Printf("Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	// if all retries failed
	return nil, err
}
