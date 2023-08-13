package main

import (
	"net/http"
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

var (
  responseTimeHistogram = prometheus.NewHistogram(
    prometheus.HistogramOpts{
      Name: "proxy_go_response_time_seconds",
      Help: "Response time in seconds",
      Buckets: prometheus.LinearBuckets(0.01, 0.1, 10),
    },
  )
)

var (
	httpStatusCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
					Name: "proxy_go_http_status_count",
					Help: "HTTP status codes count",
			},
			[]string{"code"},
	)
)

var (
	throughputHistogram = prometheus.NewHistogram(
			prometheus.HistogramOpts{
					Name:    "proxy_go_throughput_bytes",
					Help:    "Throughput in bytes",
					Buckets: prometheus.LinearBuckets(1000, 1000, 10),
			},
	)
)

func init() { 
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(responseTimeHistogram)
	prometheus.MustRegister(httpStatusCount)
	prometheus.MustRegister(throughputHistogram)
}

func startMetricsServer() *http.Server {
  http.Handle("/metrics", promhttp.Handler())
  server := &http.Server{
    Addr: ":8090",
  }
  
  return server
}