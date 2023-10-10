package metrics

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "proxy_go_requests_total",
			Help: "Total number of requests",
		},
	)
	ResponseTimeHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "proxy_go_response_time_seconds",
			Help:    "Response time in seconds",
			Buckets: prometheus.LinearBuckets(0.01, 0.1, 10),
		},
	)
	HttpStatusCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_go_http_status_count",
			Help: "HTTP status codes count",
		},
		[]string{"code"},
	)
	ThroughputHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "proxy_go_throughput_bytes",
			Help:    "Throughput in bytes",
			Buckets: prometheus.LinearBuckets(1000, 1000, 10),
		},
	)
)

func init() {
	prometheus.MustRegister(
		RequestsTotal,
		ResponseTimeHistogram,
		HttpStatusCount,
		ThroughputHistogram,
	)
}

func CreateMetricsServer() *http.Server {
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr: ":" + os.Getenv("METRICS_SERVER_PORT"),
	}

	return server
}
