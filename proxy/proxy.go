package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MaksimPozharskiy/proxy-go/metrics"
	"github.com/MaksimPozharskiy/proxy-go/server"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

var proxyTargetUrl = os.Getenv("PROXY_TARGET_URL")
var proxyServerPort = os.Getenv("PROXY_SERVER_PORT")

func RunProxyServer(ctx context.Context) {
	http.HandleFunc("/", handleRequest)

	srv := server.New(http.HandlerFunc(handleRequest), proxyServerPort)

	if err := srv.Run(ctx); err != nil {
		err = fmt.Errorf("run proxy server: %w", err)
	}
}

func handleRequest(writer http.ResponseWriter, req *http.Request) {
	startReq := time.Now()

	targetUrl := proxyTargetUrl + fmt.Sprintf("%s", req.URL)
	if targetUrl == "" {
		log.Fatal("PROXY_TARGET_URL in env not setted")
	}

	log.Println(req.Method, req.URL)

	proxyReq, err := http.NewRequest(req.Method, targetUrl, req.Body)
	if err != nil {
		http.Error(writer, "Error creaing proxy request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = req.Header.Clone()

	resp, err := getRequestWithRetry(proxyReq)

	if err != nil {
		log.Println(err)
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

	log.Printf("Request performs  %vs.\n", finishReq)
	metrics.ResponseTimeHistogram.Observe(finishReq)
	metrics.RequestsTotal.Inc()
	metrics.HttpStatusCount.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
	metrics.ThroughputHistogram.Observe(throughput)
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

		log.Printf("Request error: %v\n", err)
		log.Printf("Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}

	// if all retries failed
	return nil, err
}
