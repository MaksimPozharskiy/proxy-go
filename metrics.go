package main

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startMetricsServer() *http.Server {
  http.Handle("/metrics", promhttp.Handler())
  server := &http.Server{
    Addr: ":8090",
  }
  
  return server
}